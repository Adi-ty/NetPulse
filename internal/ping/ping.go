package ping

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Adi-ty/NetPulse/internal/config"
)

type Result struct {
	Endpoint  string
	Duration  time.Duration
	Success   bool
	Error     string
	StartTime time.Time
}

type Runner struct {
	cfg       *config.Config
	client    *http.Client
	tcpDialer *net.Dialer
}

func NewRunner(cfg *config.Config) *Runner {
	return &Runner{
		cfg: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		tcpDialer: &net.Dialer{Timeout: cfg.Timeout},
	}
}

func (r *Runner) Run(ctx context.Context, results chan<- Result) {
	var wg sync.WaitGroup

	for _, endpoint := range r.cfg.Endpoints {
		wg.Add(1)
		go func(e string) {
			defer wg.Done()
			r.pingEndpoint(ctx, e, results)
		}(endpoint)
	}

	wg.Wait()
	close(results)
}

func (r *Runner) pingEndpoint(ctx context.Context, endpoint string, results chan<- Result) {
	var wg sync.WaitGroup
	jobs := make(chan struct{}, r.cfg.Requests)

	for i := 0; i < r.cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					results <- r.executeRequest(endpoint)
				}
			}
		}()
	}

	for i := 0; i < r.cfg.Requests; i++ {
		jobs <- struct{}{}
	}
	close(jobs)
	wg.Wait()
}

func (r *Runner) executeRequest(endpoint string) Result {
	start := time.Now()
	result := Result{
		Endpoint:  endpoint,
		StartTime: start,
	}

	// Try HTTP first
	resp, err := r.client.Get(endpoint)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			result.Duration = time.Since(start)
			result.Success = true
			return result
		}
	}

	// Fallback to TCP
	conn, err := r.tcpDialer.Dial("tcp", endpoint)
	if err == nil {
		defer conn.Close()
		result.Duration = time.Since(start)
		result.Success = true
		return result
	}

	result.Error = fmt.Sprintf("connection failed: %v", err)
	return result
}