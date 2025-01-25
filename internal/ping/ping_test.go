package ping

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Adi-ty/NetPulse/internal/config"
)

func TestPingEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
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