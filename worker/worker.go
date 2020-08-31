package worker

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pb "github.com/amazingchow/photon-dance-mapreduce/api"
	"github.com/amazingchow/photon-dance-mapreduce/backend/storage"
)

type WorkerService struct {
	conf *ServiceConfig
	conn *grpc.ClientConn

	MapFun    func(string, string) []storage.KeyValue
	ReduceFun func(string, []string) string

	StopSig chan struct{}

	Storage storage.Persister
}

func NewWorkerService(conf *ServiceConfig,
	mapf func(string, string) []storage.KeyValue,
	reducef func(string, []string) string,
) *WorkerService {

	w := WorkerService{}
	w.conf = conf
	w.MapFun = mapf
	w.ReduceFun = reducef
	w.StopSig = make(chan struct{})

	opts := grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		grpc_opentracing.UnaryClientInterceptor(),
	))
	params := keepalive.ClientParameters{
		Time:                time.Second * 60,
		Timeout:             time.Second * 10,
		PermitWithoutStream: true,
	}
	conn, err := grpc.Dial(conf.MasterGRPCEndpoint, grpc.WithInsecure(), opts, grpc.WithKeepaliveParams(params))
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to connect to master node <%s>", conf.MasterGRPCEndpoint)
	}
	w.conn = conn

	var persister storage.Persister
	if conf.S3 != nil && conf.S3.Endpoint != "" {
		persister, err = storage.NewS3Persister(conf.S3)
		if err != nil {
			log.Fatal().Err(err)
		}
	} else {
		persister = storage.NewLocalFilePersister(conf.DumpRootPath)
	}
	if err = persister.Init(); err != nil {
		log.Fatal().Err(err)
	}
	w.Storage = persister

	return &w
}

// Start starts mapreduce worker node service.
func (w *WorkerService) Start() error {
	go w.run()

	log.Info().Msg("start mapreduce worker node service")
	return nil
}

// Stop stops mapreduce worker node service.
func (w *WorkerService) Stop() error {
	w.StopSig <- struct{}{}

	return nil
}

func (w *WorkerService) askForTask() *pb.IntercomResponse {
	cli := pb.NewMapReduceRPCServiceClient(w.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := pb.IntercomRequest{}
	args.MsgType = pb.IntercomType_INTERCOM_TYPE_ASK_TASK
	args.MsgContent = ""
	args.Extra = ""

	reply, err := cli.Intercom(ctx, &args)
	if err != nil {
		log.Warn().Err(err).Msg("failed to obtain task")
		return nil
	}

	return reply
}

func (w *WorkerService) sendInterFiles(tmpFile string, nReduceBucket int32) *pb.IntercomResponse {
	cli := pb.NewMapReduceRPCServiceClient(w.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := pb.IntercomRequest{}
	args.MsgType = pb.IntercomType_INTERCOM_TYPE_SEND_INTER_FILE
	args.MsgContent = tmpFile
	args.Extra = fmt.Sprintf("%d", nReduceBucket)

	reply, err := cli.Intercom(ctx, &args)
	if err != nil {
		log.Warn().Err(err).Msg("failed to send intermediate file")
		return nil
	}

	return reply
}

func (w *WorkerService) finishTask(msgType pb.IntercomType, task string) *pb.IntercomResponse {
	cli := pb.NewMapReduceRPCServiceClient(w.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	args := pb.IntercomRequest{}
	args.MsgType = msgType
	args.MsgContent = task
	args.Extra = ""

	reply, err := cli.Intercom(ctx, &args)
	if err != nil {
		log.Warn().Err(err).Msg("failed to finish task")
		return nil
	}

	return reply
}

func (w *WorkerService) runMapMode(reply *pb.IntercomResponse) {
	fd, err := os.Open(reply.GetFile())
	if err != nil {
		log.Error().Err(err).Msgf("failed to do map, since cannot open %s", reply.GetFile())
		return
	}
	defer fd.Close()

	content, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Error().Err(err).Msgf("failed to do map, since cannot read %v", reply.GetFile())
		return
	}

	kvs := w.MapFun(reply.GetFile(), string(content))
	kvp := partition(kvs, int(reply.NReduce))
	var i int32 = 0
	for ; i < reply.GetNReduce(); i++ {
		var dump string

		file := storage.IndexFile{
			Temporary: true,
			MapIdx:    reply.GetMapTaskAllocated(),
			ReduceIdx: i,
			Body:      &(kvp[i]),
		}
		if _, err = w.Storage.Writable(file); err == nil {
			if dump, err = w.Storage.Commit(file); err == nil {
				w.Storage.Abort(file) // nolint
			}
		}
		if dump != "" {
			w.sendInterFiles(dump, i)
		} else {
			log.Error().Msg("failed to do map, since something wrong happened at backend storatge")
		}
	}

	w.finishTask(pb.IntercomType_INTERCOM_TYPE_FINISH_MAP_TASK, reply.GetFile())
}

func partition(kvs []storage.KeyValue, nReduce int) [][]storage.KeyValue {
	kvp := make([][]storage.KeyValue, nReduce)
	for i := 0; i < nReduce; i++ {
		kvp[i] = make([]storage.KeyValue, 0)
	}
	for _, kv := range kvs {
		idx := fnv_1a_32(kv.Key) % uint32(nReduce)
		kvp[idx] = append(kvp[idx], kv)
	}
	return kvp
}

func (w *WorkerService) runReduceMode(reply *pb.IntercomResponse) {
	intermediates := make([]storage.KeyValue, 0)
	for _, f := range reply.GetReduceFiles() {
		mIdx, rIdx := w.Storage.RetrieveMRIdx(f)

		file := storage.IndexFile{
			Temporary: true,
			MapIdx:    mIdx,
			ReduceIdx: rIdx,
			Body:      &intermediates,
		}
		if _, err := w.Storage.Readable(file); err == nil {
			if _, err = w.Storage.Request(file); err == nil {
				w.Storage.Abort(file) // nolint
			}
		}
	}

	if len(intermediates) != 0 {
		sort.Sort(KeyValueList(intermediates))

		kvs := make([]storage.KeyValue, 0)
		i := 0
		for i < len(intermediates) {
			j := i + 1
			for j < len(intermediates) && intermediates[j].Key == intermediates[i].Key {
				j++
			}
			values := make([]string, 0)
			for k := i; k < j; k++ {
				values = append(values, intermediates[k].Value)
			}
			output := w.ReduceFun(intermediates[i].Key, values)

			kvs = append(kvs, storage.KeyValue{Key: intermediates[i].Key, Value: output})

			i = j
		}

		file := storage.IndexFile{
			Temporary: false,
			MapIdx:    -1,
			ReduceIdx: reply.GetReduceTaskAllocated(),
			Body:      &kvs,
		}
		if _, err := w.Storage.Writable(file); err == nil {
			if _, err = w.Storage.Commit(file); err == nil {
				w.Storage.Abort(file) // nolint
			}
		}

		w.finishTask(pb.IntercomType_INTERCOM_TYPE_FINISH_REDUCE_TASK, fmt.Sprintf("%d", reply.GetReduceTaskAllocated()))
	} else {
		log.Error().Msg("failed to do reduce, since something wrong happened at backend storatge")
	}
}

func (w *WorkerService) run() {
OnStopLabel:
	for {
		select {
		case <-w.StopSig:
			{
				{
					break OnStopLabel
				}
			}
		default:
			{
				reply := w.askForTask()
				if reply == nil {
					time.Sleep(time.Second)
				} else {
					switch reply.GetTaskType() {
					case pb.TaskType_TASK_TYPE_MAP:
						{
							w.runMapMode(reply)
						}
					case pb.TaskType_TASK_TYPE_REDUCE:
						{
							w.runReduceMode(reply)
						}
					}
				}
			}
		}
	}

	log.Info().Msg("stop mapreduce worker node service")
}
