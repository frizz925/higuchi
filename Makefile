GORUN=go run
GOTEST=go test -race
GOBUILD=go build -ldflags="-s -w"
GOBENCHMARK=go test -benchmem -bench=BenchmarkWorker
BUILD_OUTPUT=bin/higuchi

serve:
	$(GORUN) . serve

test:
	$(GOTEST) ./...

benchmark:
	$(GOBENCHMARK) -benchtime=10s ./internal/worker

benchmark-single-core:
	GOMAXPROCS=1 $(GOBENCHMARK) -benchtime=10s ./internal/worker

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_OUTPUT)-linux-amd64 .

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_OUTPUT)-darwin-amd64 .
