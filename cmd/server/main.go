package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/troydai/http-over-uds/internal/benchserver"
)

func main() {
	server, err := benchserver.New()
	if err != nil {
		log.Fatalln("fail to start the server", err.Error())
	}

	interuption := make(chan os.Signal, 1)
	signal.Notify(interuption, os.Interrupt)
	<-interuption
	fmt.Println("begin terminating the service")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_ = server.Shutdown(ctx)
	fmt.Println("server shut down")
}
