# NetPulse 🌐⚡

**A Concurrent Network Diagnostics Tool for HTTP/TCP Endpoints**

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

NetPulse is a powerful CLI tool for analyzing network connectivity and performance across multiple endpoints. Perfect for DevOps engineers, SREs, and network administrators who need to:

- Validate deployment connectivity
- Troubleshoot production issues
- Monitor service health
- Benchmark network performance

## Features ✨

- **Dual Protocol Support**: Tests both HTTP(S) and raw TCP connections
- **Concurrent Checks**: Simultaneous requests with configurable workers
- **Detailed Metrics**:
  - Success rate percentage
  - Latency statistics (Min/Max/Avg/P95)
  - Error aggregation and reporting
- **Docker Ready**: Pre-built image for containerized environments

## Usage 📦

```go
netpulse -n 100 -c 20 -t 2s \
  https://service1.example.com \
  service2.example.com:443 \
  192.168.1.1:3306
```

## TCP-only check

```go
netpulse -n 50 postgres-host:5432 nats-host:4222
```

## Command Flags:

| Flag | Description                     | Default |
| ---- | ------------------------------- | ------- |
| -n   | Requests per endpoint           | 5       |
| -c   | Concurrent workers per endpoint | 5       |
| -t   | Timeout per request             | 5s      |

## License 📄

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
