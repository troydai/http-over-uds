# http-over-uds

Benchmark the HTTP over UDS performance in a docker compose environment

## Usage

```
$ make benchmark
```

## Configuration

- Update `HOU_BENCHMARK_CONCURRENCY` through compose.yaml to update the concurrency settings of the runs
- Update `HOU_BENCHMARK_DURATION` through compose.yaml to update the duration of the run

The CPU limits on the server is set to 0.1 to constrain the server capacity

