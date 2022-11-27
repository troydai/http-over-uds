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

## Sample result

```
http-over-uds-client-1  |                     Count  Error     p99     p95     p50    Status
http-over-uds-client-1  |    Total 1 Clients   1602      0    95.4    89.4     0.5  200=100%
http-over-uds-client-1  |    Total 4 Clients   2835      0    99.8    95.2     0.9  200=100%
http-over-uds-client-1  |   Total 16 Clients   3970      0   190.8   101.3     3.2  200=100%
http-over-uds-client-1  |   Total 64 Clients   3230      0   800.1   498.5   193.4  200=100%
http-over-uds-client-1  |  Total 256 Clients   4634      0   998.9   896.4   599.1  200=100%
http-over-uds-client-1  | Total 1024 Clients   3712      0  6499.2  5901.0  3695.2  200=100%
```