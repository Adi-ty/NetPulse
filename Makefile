BINARY = netpulse
IMAGE_TAG = v1.0.0
GO_TEST_FLAGS = -race -covermode=atomic

.PHONY: build-docker test

build:
	go build -o bin/$(BINARY) ./cmd/netpulse

build-docker:
	docker build --platform linux/amd64 -t $(BINARY):$(IMAGE_TAG) .

run-docker:
	docker run --rm $(BINARY):$(IMAGE_TAG) -n 5 google.com:443

test:
	go test $(GO_TEST_FLAGS) -coverprofile=coverage.txt ./...

lint:
	golangci-lint run