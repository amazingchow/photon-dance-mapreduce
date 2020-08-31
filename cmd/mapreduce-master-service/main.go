package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/amazingchow/photon-dance-mapreduce/api"
	"github.com/amazingchow/photon-dance-mapreduce/master"
	"github.com/amazingchow/photon-dance-mapreduce/utils"
)

var (
	_ConfigPath = flag.String("conf", "conf/master_conf.json", "worker config")
	_Level      = flag.String("level", "info", "log level, options [debug info warn error]")
)

type MapReduceIngressServer struct {
	executor *master.MasterService
}

func newMapReduceIngressServer(conf *master.ServiceConfig) *MapReduceIngressServer {
	return &MapReduceIngressServer{
		executor: master.NewMasterService(conf),
	}
}

func (mris *MapReduceIngressServer) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	if len(req.GetTask().GetInputs()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty input")
	}

	if err := mris.executor.AddTask(ctx, req.GetTask()); err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	return &pb.AddTaskResponse{}, nil
}

func (mris *MapReduceIngressServer) ListWorkers(ctx context.Context, req *pb.ListWorkersRequest) (*pb.ListWorkersResponse, error) {
	return &pb.ListWorkersResponse{}, nil
}

func (mris *MapReduceIngressServer) Intercom(ctx context.Context, req *pb.IntercomRequest) (*pb.IntercomResponse, error) {
	reply := pb.IntercomResponse{}

	if err := mris.executor.InterComm(ctx, req, &reply); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &reply, nil
}

func serverGrpcService(ctx context.Context, mris *MapReduceIngressServer, conf *master.ServiceConfig, stopGroup *sync.WaitGroup, stopCh chan struct{}) {
	stopGroup.Add(1)
	defer stopGroup.Done()

	lis, err := net.Listen("tcp", conf.GRPCEndpoint)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start grpc service")
	}

	opts := []grpc.ServerOption{
		grpc.MaxSendMsgSize(64 * 1024 * 1024),
		grpc.MaxRecvMsgSize(64 * 1024 * 1024),
	}
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterMapReduceRPCServiceServer(grpcServer, mris)
	log.Info().Msgf("grpc service is listening at \x1b[1;31m%s\x1b[0m", conf.GRPCEndpoint)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Warn().Err(err)
		}
	}()

GRPC_LOOP:
	for { // nolint
		select {
		case _, ok := <-stopCh:
			{
				if !ok {
					break GRPC_LOOP
				}
			}
		}
	}

	grpcServer.GracefulStop()
	log.Info().Msg("stop grpc service")
}

func sereveHTTPService(ctx context.Context, mris *MapReduceIngressServer, conf *master.ServiceConfig, stopGroup *sync.WaitGroup, stopCh chan struct{}) {
	stopGroup.Add(1)
	defer stopGroup.Done()

	mux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	if err := pb.RegisterMapReduceRPCServiceHandlerFromEndpoint(ctx, mux, conf.GRPCEndpoint, dialOpts); err != nil {
		log.Fatal().Err(err).Msg("failed to register http service")
	}

	http.Handle("/", mux)
	httpServer := http.Server{
		Addr: conf.HTTPEndpoint,
	}

	log.Info().Msgf("http service is listening at \x1b[1;31m%s\x1b[0m", conf.HTTPEndpoint)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Warn().Err(err)
		}
	}()

HTTP_LOOP:
	for { // nolint
		select {
		case _, ok := <-stopCh:
			{
				if !ok {
					break HTTP_LOOP
				}
			}
		}
	}

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Warn().Err(err).Msg("failed to stop http service")
	}
	log.Info().Msg("stop http service")
}

func main() {
	flag.Parse()

	// set logger level
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	switch *_Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// load master config
	var conf master.ServiceConfig
	if err := utils.LoadConfig(*_ConfigPath, &conf); err != nil {
		log.Fatal().Err(err).Msgf("failed to load config <%s>", *_ConfigPath)
	}

	stopGroup := &sync.WaitGroup{}
	defer func() {
		stopGroup.Wait()
	}()
	stopCh := make(chan struct{})

	mris := newMapReduceIngressServer(&conf)
	mris.executor.Start() // nolint

	// serve grpc service && http service
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()

	go serverGrpcService(ctx, mris, &conf, stopGroup, stopCh)
	go sereveHTTPService(ctx, mris, &conf, stopGroup, stopCh)

	// wait for exit signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
	// send stop signal to grpc service && http service
	close(stopCh)

	mris.executor.Stop() // nolint
}
