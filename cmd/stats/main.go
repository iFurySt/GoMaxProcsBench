package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	mode  int
	ts    time.Duration
	times int
)

func init() {
	flag.IntVar(&mode, "mode", 0, "0: auto, 1: runtime")
	flag.DurationVar(&ts, "ts", 0, "time to run")
	flag.IntVar(&times, "times", 1, "times to run")
	flag.Parse()
}

func main() {
	var total int64 = 0
	cnt := 0
	for range times {
		cmd := exec.Command("go", "run", "cmd/bench/main.go",
			"--silent", "-mode", fmt.Sprint(mode), "-ts", fmt.Sprint(ts.String()))

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Failed to execute command: %v\n", err)
			continue
		}
		qps, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
		if err != nil {
			log.Printf("Failed to parse output: %v\n", err)
			continue
		}
		total += qps
		cnt++
	}
	if cnt > 0 {
		log.Printf("Average QPS: %d\n", total/int64(cnt))
	} else {
		log.Printf("No valid result\n")
	}
}
