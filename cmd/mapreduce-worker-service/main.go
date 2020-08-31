package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/amazingchow/photon-dance-mapreduce/utils"
	"github.com/amazingchow/photon-dance-mapreduce/worker"
)

var (
	_ConfigPath = flag.String("conf", "conf/worker_1_conf.json", "worker config")
	_Level      = flag.String("level", "info", "log level, options [debug info warn error]")
	_CPUProfile = flag.String("cpuprofile", "cpu.prof", "dump cpu profile")
	_MemProfile = flag.String("memprofile", "mem.prof", "dump memory profile")
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

	{
		if *_CPUProfile != "" {
			fd, err := os.Create(*_CPUProfile)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create CPU profile file")
			}
			defer fd.Close()

			if err := pprof.StartCPUProfile(fd); err != nil {
				log.Fatal().Err(err).Msg("failed to start CPU profile")
			}
			defer pprof.StopCPUProfile()
		}

		defer func() {
			if *_MemProfile != "" {
				fd, err := os.Create(*_MemProfile)
				if err != nil {
					log.Fatal().Err(err).Msg("failed to create Mem profile file")
				}
				defer fd.Close()

				// get up-to-date statistics
				runtime.GC()
				if err := pprof.WriteHeapProfile(fd); err != nil {
					log.Fatal().Err(err).Msg("failed to start Mem profile")
				}
			}
		}()
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
