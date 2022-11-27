package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/troydai/http-over-uds/internal/benchclient"
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
		tablify.Print(os.Stdout, series, true, true)

		mergedSeries = append(mergedSeries, merged)
	}

	if len(mergedSeries) > 0 {
		fmt.Println("Summary of all the runs")
		tablify.Print(os.Stdout, mergedSeries, false, false)
	}
}

func benchmark(concurrency int, path string, duration time.Duration) []*summary.Series {
	agents := make([]*benchclient.Agent, 0, concurrency)
	for i := 0; i < int(concurrency); i++ {
		agents = append(agents, benchclient.New(path, 0))
	}

	wg := &sync.WaitGroup{}
	wg.Add(concurrency)

	dataSeries := make([]*summary.Series, len(agents))
	for idx := range dataSeries {
		dataSeries[idx] = summary.NewSeries(fmt.Sprintf("Client-%d", idx))
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	for i := 0; i < int(concurrency); i++ {
		agents[i].Start(ctx, wg, dataSeries[i].Append)
	}

	wg.Wait()

	return dataSeries
}
