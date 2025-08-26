package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
    intervalStr := os.Getenv("INTERVAL_SECONDS")
    intervalSeconds, err := strconv.Atoi(intervalStr)
    if err != nil || intervalSeconds < 1 {
        fmt.Println("Invalid interval, setting to default (30 seconds)")
        intervalSeconds = 30
    }
    ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
    defer ticker.Stop()
    for t := range ticker.C {
        fmt.Printf("Tick at %v\n", t)
    }
}
