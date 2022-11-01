package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/VladimirBlinov/TransactionService/Backend/internal/app/workerserver"
	"github.com/sirupsen/logrus"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	logger := logrus.New()

	config := workerserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		logger.Fatalf("Failed load config", err.Error())
	}

	wrkserver, err := workerserver.NewWorkerServer(config)
	if err != nil {
		logger.Fatalf("Failed create WorkerServer", err.Error())
	}

	go func() {
		if err = wrkserver.Start(); err != nil {
			logger.Fatalf("error on worker server start", err.Error())
		}
	}()

	logger.Println("Started worker server...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Println("Shutting down worker server...")
	if err = wrkserver.ShutDown(context.Background()); err != nil {
		logger.Fatalf("error on server shutdown: %s", err.Error())
	}

	os.Exit(0)
}
