package master

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	pb "github.com/amazingchow/mapreduce/api"
)

type MasterService struct {
	mu   sync.RWMutex
	conf *ServiceConfig

	MapTaskTable           map[string]pb.TaskStatus
	ReduceTaskTable        map[string]pb.TaskStatus
	MapTaskCh              chan string
	ReduceTaskCh           chan string
	NMap                   int32
	NReduce                int32
	ComputingJobArrivedSig chan struct{}
	AllMapTasksFinished    bool
	AllReduceTasksFinished bool
	AllTasksFinishedSig    chan struct{}
	StopSig                chan struct{}

	InterFiles [][]string
}

func NewMasterService(conf *ServiceConfig) *MasterService {
	m := MasterService{}
	m.conf = conf
	m.MapTaskTable = make(map[string]pb.TaskStatus)
	m.ReduceTaskTable = make(map[string]pb.TaskStatus)
	m.MapTaskCh = make(chan string)
	m.ReduceTaskCh = make(chan string)
	m.NMap = 0
	m.NReduce = conf.NReduce
	m.ComputingJobArrivedSig = make(chan struct{})
	m.AllMapTasksFinished = false
	m.AllReduceTasksFinished = false
	m.AllTasksFinishedSig = make(chan struct{})
	m.StopSig = make(chan struct{})
	m.InterFiles = make([][]string, conf.NReduce)

	return &m
}

// Start starts mapreduce master node service.
func (m *MasterService) Start() error {
	go m.run()

	log.Info().Msg("start mapreduce master node service")
	return nil
}

// Stop stops mapreduce master node service.
func (m *MasterService) Stop() error {
	m.StopSig <- struct{}{}

	return nil
}

// AddTask adds computing job to mapreduce backend service.
func (m *MasterService) AddTask(ctx context.Context, task *pb.Task) error {
	select {
	case <-m.AllTasksFinishedSig:
		{
			goto OnAddTask
		}
	default:
		{
			// if former computing job has not been finished, we should not provide service now.
			return errors.New("backend busy")
		}
	}

OnAddTask:
	log.Info().Msg("new computing job arrived")

	for task := range m.MapTaskTable {
		delete(m.MapTaskTable, task)
	}
	for _, file := range task.Inputs {
		m.MapTaskTable[file] = pb.TaskStatus_TASK_STATUS_UNALLOTED
	}

	for task := range m.ReduceTaskTable {
		delete(m.ReduceTaskTable, task)
	}
	var i int32 = 0
	for ; i < m.NReduce; i++ {
		m.ReduceTaskTable[fmt.Sprintf("%d", i)] = pb.TaskStatus_TASK_STATUS_UNALLOTED
	}

	m.InterFiles = m.InterFiles[:0]
	for i = 0; i < m.NReduce; i++ {
		m.InterFiles[i] = make([]string, 0)
	}

	m.ComputingJobArrivedSig <- struct{}{}

	return nil
}

// InterComm used for internal-communication called by workers.
func (m *MasterService) InterComm(ctx context.Context, args *pb.IntercomRequest, reply *pb.IntercomResponse) error {
	switch args.GetMsgType() {
	case pb.IntercomType_INTERCOM_TYPE_ASK_TASK:
		{
		OnExitLabel:
			for {
				select {
				case task := <-m.MapTaskCh:
					{
						reply.TaskType = pb.TaskType_TASK_TYPE_MAP
						reply.File = task
						reply.NReduce = m.NReduce
						reply.MapTaskAllocated = m.NMap

						m.mu.Lock()
						m.MapTaskTable[task] = pb.TaskStatus_TASK_STATUS_ALLOTED
						m.NMap++
						m.mu.Unlock()

						go m.takeFaultTolerantPolicy(pb.TaskType_TASK_TYPE_MAP, task)
						break OnExitLabel
					}
				case task := <-m.ReduceTaskCh:
					{
						reply.TaskType = pb.TaskType_TASK_TYPE_REDUCE
						reply.NReduce = m.NReduce
						x, _ := strconv.Atoi(task)
						reply.ReduceTaskAllocated = int32(x)
						reply.ReduceFiles = m.InterFiles[x]

						m.mu.Lock()
						m.ReduceTaskTable[task] = pb.TaskStatus_TASK_STATUS_ALLOTED
						m.mu.Unlock()

						go m.takeFaultTolerantPolicy(pb.TaskType_TASK_TYPE_REDUCE, task)
						break OnExitLabel
					}
				case <-m.AllTasksFinishedSig:
					{
						return errors.New("no tasks available")
					}
				default:
					{
						if IsContextDone(ctx) {

							return errors.New("ctx done")
						}
					}
				}
			}
		}
	case pb.IntercomType_INTERCOM_TYPE_FINISH_MAP_TASK:
		{
			m.mu.Lock()
			task := args.MsgContent
			if m.MapTaskTable[task] == pb.TaskStatus_TASK_STATUS_ALLOTED {
				m.MapTaskTable[task] = pb.TaskStatus_TASK_STATUS_DONE
			}
			m.mu.Unlock()
		}
	case pb.IntercomType_INTERCOM_TYPE_FINISH_REDUCE_TASK:
		{
			m.mu.Lock()
			task := args.MsgContent
			if m.ReduceTaskTable[task] == pb.TaskStatus_TASK_STATUS_ALLOTED {
				m.ReduceTaskTable[task] = pb.TaskStatus_TASK_STATUS_DONE
			}
			m.mu.Unlock()
		}
	case pb.IntercomType_INTERCOM_TYPE_SEND_INTER_FILE:
		{
			interFile := args.MsgContent
			nReduce, _ := strconv.Atoi(args.Extra)
			m.InterFiles[nReduce] = append(m.InterFiles[nReduce], interFile)
		}
	}

	return nil
}

func (m *MasterService) takeFaultTolerantPolicy(taskType pb.TaskType, task string) {
	timer := time.NewTicker(10 * time.Second)
	defer timer.Stop()

OnExitLabel:
	for {
		select {
		case <-timer.C:
			{
				if taskType == pb.TaskType_TASK_TYPE_MAP {
					m.mu.Lock()
					m.MapTaskTable[task] = pb.TaskStatus_TASK_STATUS_UNALLOTED
					m.mu.Unlock()

					log.Debug().Msgf("map task <%s> timeouts, reschedule it again", task)
					m.MapTaskCh <- task
				} else if taskType == pb.TaskType_TASK_TYPE_REDUCE {
					m.mu.Lock()
					m.ReduceTaskTable[task] = pb.TaskStatus_TASK_STATUS_UNALLOTED
					m.mu.Unlock()

					log.Debug().Msgf("reduce task <%s> timeouts, reschedule it again", task)
					m.ReduceTaskCh <- task
				}
				break OnExitLabel
			}
		default:
			{
				if taskType == pb.TaskType_TASK_TYPE_MAP {
					m.mu.RLock()
					if m.MapTaskTable[task] == pb.TaskStatus_TASK_STATUS_DONE {
						m.mu.RUnlock()
						log.Debug().Msgf("map task <%s> done", task)
						break OnExitLabel
					}
					m.mu.RUnlock()
				} else if taskType == pb.TaskType_TASK_TYPE_REDUCE {
					m.mu.RLock()
					if m.ReduceTaskTable[task] == pb.TaskStatus_TASK_STATUS_DONE {
						m.mu.RUnlock()
						log.Debug().Msgf("reduce task <%s> done", task)
						break OnExitLabel
					}
					m.mu.RUnlock()
				}
			}
		}
	}
}

func (m *MasterService) waitForAllMapTasks() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, status := range m.MapTaskTable {
		if status != pb.TaskStatus_TASK_STATUS_DONE {
			return false
		}
	}

	return true
}

func (m *MasterService) waitForAllReduceTasks() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, status := range m.ReduceTaskTable {
		if status != pb.TaskStatus_TASK_STATUS_DONE {
			return false
		}
	}

	return true
}

func (m *MasterService) run() {
OnStopLabel:
	for {
		select {
		case <-m.ComputingJobArrivedSig:
			{
				for task := range m.MapTaskTable {
					m.MapTaskCh <- task
				}
				for !m.waitForAllMapTasks() {
					time.Sleep(200 * time.Millisecond)
				}
				m.AllMapTasksFinished = true

				for task := range m.ReduceTaskTable {
					m.ReduceTaskCh <- task
				}
				for !m.waitForAllReduceTasks() {
					time.Sleep(200 * time.Millisecond)
				}
				m.AllReduceTasksFinished = true

				m.AllTasksFinishedSig <- struct{}{}

				log.Info().Msg("all tasks done")
			}
		case <-m.StopSig:
			{
				break OnStopLabel
			}
		}
	}

	log.Info().Msg("stop mapreduce master node service")
}
