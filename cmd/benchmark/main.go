package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/troydai/http-over-uds/internal/summary"
	"github.com/troydai/http-over-uds/internal/tablify"
)

type envVars struct {
	UDSPath     string        `env:"HOU_BENCHMARK_UDS"`
	Concurrency []int         `env:"HOU_BENCHMARK_CONCURRENCY" envSeparator:","`
	Duration    time.Duration `env:"HOU_BENCHMARK_DURATION"`
}

func main() {
	var ev envVars
	if err := env.Parse(&ev); err != nil {
		log.Fatalf("fail to load environment variables: %v", err.Error())
	}

	fmt.Println("Starting warming run")
	benchmark(1, ev.UDSPath, 5*time.Second)

	var mergedSeries []*summary.Series
	for _, c := range ev.Concurrency {
		fmt.Printf("Starting run with %d concurrency.\n", c)
		series := benchmark(c, ev.UDSPath, ev.Duration)
		merged := summary.Merge(fmt.Sprintf("Total %d Clients", c), series...)

		series = append(series, merged)
		tablify.Print(os.Stdout, series)

		mergedSeries = append(mergedSeries, merged)
	}

	if len(mergedSeries) > 0 {
		fmt.Println("Summary of all the runs")
		tablify.Print(os.Stdout, mergedSeries)
	}
}

func benchmark(concurrency int, path string, duration time.Duration) []*summary.Series {
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
