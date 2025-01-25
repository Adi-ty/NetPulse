package config

import (
	"flag"
	"time"
)

type Config struct {
	Requests    int
	Concurrency int
	Timeout     time.Duration
	Endpoints   []string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ParseFlags() {
	flag.IntVar(&c.Requests, "n", 5, "Number of requests to perform")
	flag.IntVar(&c.Concurrency, "c", 5, "Number of concurrent requests")
	flag.DurationVar(&c.Timeout, "t", 5*time.Second, "Timeout for each request")
	flag.Parse()
	c.Endpoints = flag.Args()
}

func (c *Config) TotalRequests() int {
	return c.Requests * len(c.Endpoints)
}