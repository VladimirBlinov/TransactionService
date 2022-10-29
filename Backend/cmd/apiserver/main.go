package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/BurntSushi/toml"
	"github.com/VladimirBlinov/MarketPlace/MarketPlace/internal/app/apiserver"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	logger := logrus.New()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		logger.Fatal(err)
	}

	apisrv := new(apiserver.ApiServer)

	go func() {
		if err = apisrv.Start(config); err != nil {
			if err != http.ErrServerClosed {
				logger.Fatalf("error on server start: %s", err.Error())
			}
		}
	}()

	logger.Println("Started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Println("Shutting down...")
	if err = apisrv.ShutDown(context.Background()); err != nil {
		logger.Fatalf("error on server shutdown: %s", err.Error())
	}
	os.Exit(0)
}
