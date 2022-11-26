package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
)

const _udsPath = "socket/server.socks"

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("fail to create a logger: %s", err.Error())
	}

	listener, err := net.Listen("unix", "socket/server.socks")
	if err != nil {
		logger.Fatal("fail to create TCP listner", zap.Error(err))
	}

	if _, err := os.Stat(_udsPath); os.IsExist(err) {
		if err := os.Remove(_udsPath); err != nil {
			logger.Fatal("fail to remove pre-exist socket path")
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pass: http over uds"))
	})

	server := http.Server{
		Handler: mux,
	}

	go func() {
		logger.Info("start serving traffic", zap.String("address", _udsPath))
		server.Serve(listener)
	}()

	interuption := make(chan os.Signal, 1)
	signal.Notify(interuption, os.Interrupt)
	<-interuption
	logger.Info("begin terminating the service")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = server.Shutdown(ctx)
	logger.Info("server shut down", zap.Error(err))
}
