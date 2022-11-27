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
	"github.com/troydai/http-over-uds/internal/tablify"
)

type envVars struct {
	UDSPath     string        `env:"HOU_BENCHMARK_UDS"`
	Concurrency []int         `env:"HOU_BENCHMARK_CONCURRENCY" envSeparator:","`
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

	if len(ev.Concurrency) == 1 {
		series := benchmark(ev.Concurrency[0], ev.UDSPath, ev.Duration, false)
		series = append(series, summary.Merge("Total", series...))
		for _, l := range tablify.GetLines(series) {
			fmt.Println(l)
		}
	} else {
		var series []*summary.Series
		for _, c := range ev.Concurrency {
			series = append(series, benchmark(c, ev.UDSPath, ev.Duration, true)...)
		}
		for _, l := range tablify.GetLines(series) {
			fmt.Println(l)
		}
	}
}

func benchmark(concurrency int, path string, duration time.Duration, summaryOnly bool) []*summary.Series {
	clients := make([]*http.Client, 0, concurrency)
	for i := 0; i < int(concurrency); i++ {
		client := http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					var d net.Dialer
					return d.DialContext(ctx, "unix", path)
				},
			},
		}
		clients = append(clients, &client)
	}

	wg := &sync.WaitGroup{}
	wg.Add(concurrency)

	dataSeries := make([]*summary.Series, len(clients))
	for idx := range dataSeries {
		dataSeries[idx] = summary.NewSeries(fmt.Sprintf("Client-%d", idx))
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	for i := 0; i < int(concurrency); i++ {
		startClient(ctx, wg, clients[i], dataSeries[i])
	}

	wg.Wait()

	if summaryOnly {
		return []*summary.Series{summary.Merge(fmt.Sprintf("%3d Clients", concurrency), dataSeries...)}
	}

	return dataSeries
}

func startClient(ctx context.Context, wg *sync.WaitGroup, client *http.Client, series *summary.Series) {
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
