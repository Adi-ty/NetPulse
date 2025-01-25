package stats

import (
	"fmt"
	"sort"
	"time"

	"github.com/fatih/color"

	"github.com/Adi-ty/NetPulse/internal/ping"
)

type Processor struct {
	endpoints map[string]*EndpointStats
}

type EndpointStats struct {
	Durations []float64
	Errors    []string
}

func NewProcessor(endpoints []string) *Processor {
	p := &Processor{
		endpoints: make(map[string]*EndpointStats),
	}
	
	for _, ep := range endpoints {
		p.endpoints[ep] = &EndpointStats{}
	}
	return p
}

func (p *Processor) Process(results <-chan ping.Result) {
	for result := range results {
		stats := p.endpoints[result.Endpoint]
		if result.Success {
			stats.Durations = append(stats.Durations, result.Duration.Seconds()*1000)
		} else {
			stats.Errors = append(stats.Errors, result.Error)
		}
	}
}

func (p *Processor) PrintSummary(duration time.Duration) {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Printf("\n%s\n", cyan("=== NetPulse Summary ==="))
	fmt.Printf("Total time: %s\n", yellow(duration.Round(time.Millisecond)))

	for endpoint, stats := range p.endpoints {
		total := len(stats.Durations) + len(stats.Errors)
		successRate := float64(len(stats.Durations)) / float64(total) * 100
		sort.Float64s(stats.Durations)

		fmt.Printf("\n%s %s\n", cyan("Endpoint:"), green(endpoint))
		fmt.Printf("  Success rate: %.1f%%\n", successRate)
		fmt.Printf("  Requests: %d successful, %d failed\n", len(stats.Durations), len(stats.Errors))
		
		if len(stats.Durations) > 0 {
			fmt.Printf("  Latency (ms):\n")
			fmt.Printf("    Min: %.2f\n", stats.Durations[0])
			fmt.Printf("    Max: %.2f\n", stats.Durations[len(stats.Durations)-1])
			fmt.Printf("    Avg: %.2f\n", average(stats.Durations))
			fmt.Printf("    P95: %.2f\n", percentile(stats.Durations, 95))
		}
	}

	p.printErrors()
}

func (p *Processor) printErrors() {
	hasErrors := false
	for _, stats := range p.endpoints {
		if len(stats.Errors) > 0 {
			hasErrors = true
			break
		}
	}

	if !hasErrors {
		return
	}

	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("\n%s\n", red("=== Errors ==="))
	for endpoint, stats := range p.endpoints {
		if len(stats.Errors) > 0 {
			fmt.Printf("%s:\n", endpoint)
			for _, e := range stats.Errors {
				fmt.Printf("  - %s\n", e)
			}
		}
	}
}

func average(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func percentile(data []float64, p float64) float64 {
	if len(data) == 0 {
		return 0
	}
	index := int(float64(len(data)-1) * p / 100)
	return data[index]
}