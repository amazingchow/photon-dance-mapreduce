package main

import (
	"context"
	"flag"
	"net"
	"net/http"
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

	pbapi "github.com/amazingchow/mapreduce/api"
	"github.com/amazingchow/mapreduce/master"
	"github.com/amazingchow/mapreduce/utils"
)

var (
	_ConfigPath = flag.String("conf", "config/config.json", "worker config")
	_Level      = flag.String("level", "info", "log level, options [debug info warn error]")
)

type MapReduceIngressServer struct {
}

func newMapReduceIngressServer() *MapReduceIngressServer {
	return &MapReduceIngressServer{}
}

func (mris *MapReduceIngressServer) AddTask(ctx context.Context, req *pbapi.AddTaskRequest) (*pbapi.AddTaskResponse, error) {
	if len(req.GetTask().GetInputs()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty input")
	}

	return &pbapi.AddTaskResponse{}, nil
}

func (mris *MapReduceIngressServer) ListWorkers(ctx context.Context, req *pbapi.ListWorkersRequest) (*pbapi.ListWorkersResponse, error) {
	return &pbapi.ListWorkersResponse{}, nil
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

	pbapi.RegisterManagerServiceServer(grpcServer, mris)
	log.Info().Msgf("grpc is listening at %s", conf.GRPCEndpoint)
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
		default:
			{

			}
		}
	}

	grpcServer.GracefulStop()
	log.Info().Msg("stop grpc service")
}

func sereveHTTPService(ctx context.Context, mris *MapReduceIngressServer, config *master.ServiceConfig, stopGroup *sync.WaitGroup, stopCh chan struct{}) {
	stopGroup.Add(1)
	defer stopGroup.Done()

	mux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	if err := pbapi.RegisterManagerServiceHandlerFromEndpoint(ctx, mux, config.GRPCEndpoint, dialOpts); err != nil {
		log.Fatal().Err(err).Msg("failed to register grpc gateway")
	}

	http.Handle("/", mux)
	httpServer := http.Server{
		Addr: config.HTTPEndpoint,
	}

	log.Info().Msgf("grpc gateway is listening at %s", config.HTTPEndpoint)
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
		default:
			{

			}
		}
	}

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Warn().Err(err).Msg("failed to stop grpc gateway")
	}
	log.Info().Msg("stop grpc gateway")
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
		log.Info().Msg("stop mapreduce master service")
	}()
	stopCh := make(chan struct{})

	mris := newMapReduceIngressServer()

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
MAIN_LOOP:
	for { // nolint
		select {
		case <-sigCh:
			{
				// send stop signal to grpc service && http service
				close(stopCh)
				break MAIN_LOOP
			}
		default:
			{

			}
		}
	}
}
