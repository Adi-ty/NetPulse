package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"

	"github.com/Adi-ty/NetPulse/internal/config"
	"github.com/Adi-ty/NetPulse/internal/ping"
	"github.com/Adi-ty/NetPulse/internal/stats"
)

func main() {
	cfg := config.NewConfig()
	cfg.ParseFlags()

	if len(cfg.Endpoints) == 0 {
		color.Red("No endpoints provided")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupSignalHandler(cancel)

	runner := ping.NewRunner(cfg)
	results := make(chan ping.Result, cfg.TotalRequests())
	
	start := time.Now()
	go runner.Run(ctx, results)

	processor := stats.NewProcessor(cfg.Endpoints)
	processor.Process(results)

	processor.PrintSummary(time.Since(start))
}

func setupSignalHandler(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()
}