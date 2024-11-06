package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/automaxprocs/maxprocs"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
)

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

var (
	mode   int
	ts     time.Duration
	silent bool
)

func init() {
	flag.IntVar(&mode, "mode", 0, "0: auto, 1: runtime")
	flag.DurationVar(&ts, "ts", 0, "time to run")
	flag.BoolVar(&silent, "silent", false, "silent mode")
	flag.Parse()
}

func main() {
	Printf("mode: %d, ts: %s\n", mode, ts)

	if mode == 1 {
		runtime.GOMAXPROCS(runtime.NumCPU())
		Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	} else {
		_, _ = maxprocs.Set(maxprocs.Logger(Printf))
	}

	var (
		st     = time.Now()
		count  atomic.Int64
		sigs   = make(chan os.Signal, 1)
		ctx    context.Context
		cancel context.CancelFunc
	)

	if ts > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), ts)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer func() {
		Printf("count: %d, time: %v, qps: %.0f\n", count.Load(), time.Since(st),
			float64(count.Load())/time.Since(st).Seconds())
		if silent {
			fmt.Printf("%.0f\n", float64(count.Load())/time.Since(st).Seconds())
		}
	}()

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return
		case <-sigs:
			return
		default:
			go func() {
				_ = fib(10)
				count.Add(1)
			}()
		}
	}
}

func Printf(format string, v ...interface{}) {
	if silent {
		return
	}
	log.Printf(format, v...)
}
