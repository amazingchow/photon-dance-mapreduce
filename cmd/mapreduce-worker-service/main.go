package main

import (
	"flag"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/amazingchow/mapreduce/utils"
	"github.com/amazingchow/mapreduce/worker"
)

var (
	_ConfigPath = flag.String("conf", "conf/worker_1_conf.json", "worker config")
	_Level      = flag.String("level", "info", "log level, options [debug info warn error]")
)

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

	// load worker config
	var conf worker.ServiceConfig
	if err := utils.LoadConfig(*_ConfigPath, &conf); err != nil {
		log.Fatal().Err(err).Msgf("failed to load config <%s>", *_ConfigPath)
	}

	mapf, reducef := utils.LoadPlugin()
	executor := worker.NewWorkerService(&conf, mapf, reducef)
	executor.Start() // nolint

	// wait for exit signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh

	executor.Stop() // nolint
}
