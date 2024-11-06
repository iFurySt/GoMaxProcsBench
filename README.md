# GoMaxProcsBench

Compare the performance of `runtime.GOMAXPROCS` in the case that may be affected by Cgroup environment.

## Quick Start

```shell
# 1. start the container
docker compose up -d
docker compose exec golang bash

# 2. run the bench in different mode
go run cmd/bench/main.go --ts 5s --mode 0
go run cmd/bench/main.go --ts 5s --mode 1

# 3. run the stats to get the result in different mode
go run cmd/stats/main.go --ts 5s --mode 0 --times 2
go run cmd/stats/main.go --ts 5s --mode 1 --times 2
```
