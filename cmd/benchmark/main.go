package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"github.com/troydai/http-over-uds/internal/summary"
)

type envVars struct {
	UdsPath     string        `env:"HOU_BENCHMARK_UDS"`
	Concurrency int           `env:"HOU_BENCHMARK_CONCURRENCY"`
	Duration    time.Duration `env:"HOU_BENCHMARK_DURATION"`
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("fail to create logger: %s", err.Error())
	}

	var ev envVars
	if err := env.Parse(&ev); err != nil {
		logger.Fatal("fail to load environment variables", zap.Error(err))
	}

	clients := make([]*http.Client, 0, ev.Concurrency)
	for i := 0; i < int(ev.Concurrency); i++ {
		client := http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					var d net.Dialer
					return d.DialContext(ctx, "unix", ev.UdsPath)
				},
			},
		}
		clients = append(clients, &client)
	}

	wg := &sync.WaitGroup{}
	wg.Add(ev.Concurrency)

	dataSeries := make([]*summary.Series, len(clients))
	for idx := range dataSeries {
		dataSeries[idx] = summary.NewSeries()
	}

	ctx, cancel := context.WithTimeout(context.Background(), ev.Duration)
	defer cancel()

	for i := 0; i < int(ev.Concurrency); i++ {
		startClient(ctx, wg, clients[i], dataSeries[i], logger)
	}

	wg.Wait()

	for i, s := range dataSeries {
		logger.Info(s.Summary(fmt.Sprintf("Client %02d", i)))
	}

	merged := summary.Merge(dataSeries...)
	logger.Info(merged.Summary("    Total"))
}

func startClient(ctx context.Context, wg *sync.WaitGroup, client *http.Client, series *summary.Series, logger *zap.Logger) {
	go func() {
		defer wg.Done()
		for {
			if ctx.Err() != nil {
				return
			}

			beginning := time.Now()
			resp, err := client.Get("http://localhost")
			series.Append(resp, err, beginning)
		}
	}()
}
