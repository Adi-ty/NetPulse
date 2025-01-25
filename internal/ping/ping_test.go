package ping

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Adi-ty/NetPulse/internal/config"
)

func TestPingEndpoint(t *testing.T) {
    // Test HTTP Server
    httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))
    defer httpServer.Close()

    // Test TCP Server
    tcpListener, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        t.Fatalf("Failed to start TCP server: %v", err)
    }
    defer tcpListener.Close()
    
    go func() {
        for {
            conn, err := tcpListener.Accept()
            if err != nil {
                return
            }
            conn.Close()
        }
    }()

    cfg := &config.Config{
        Requests:    5,
        Concurrency: 2,
        Timeout:     1 * time.Second,
        Endpoints: []string{
            httpServer.URL,          // HTTP test endpoint
            tcpListener.Addr().String(), // TCP test endpoint
        },
    }

    runner := NewRunner(cfg)
    results := make(chan Result, 10)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    go runner.Run(ctx, results)

    successCount := 0
    for res := range results {
        if res.Success {
            successCount++
        }
        if res.Endpoint == httpServer.URL && res.Duration == 0 {
            t.Error("HTTP request duration not recorded")
        }
    }

    expectedSuccess := cfg.Requests * len(cfg.Endpoints)
    if successCount != expectedSuccess {
        t.Errorf("Expected %d successes, got %d", expectedSuccess, successCount)
    }
}