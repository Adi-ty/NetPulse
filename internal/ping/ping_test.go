package ping

import (
	"context"
	"testing"
	"time"

	"github.com/Adi-ty/NetPulse/internal/config"
)

func TestPingEndpoint(t *testing.T) {
	cfg := &config.Config{
		Requests:    3,
		Concurrency: 2,
		Timeout:     1 * time.Second,
	}

	runner := NewRunner(cfg)
	results := make(chan Result, 3)

	go func() {
		runner.Run(context.Background(), results)
		close(results)
	}()

	successCount := 0
	for res := range results {
		if res.Success {
			successCount++
		}
	}

	if successCount == 0 {
		t.Fatal("No successful pings to localhost")
	}
}