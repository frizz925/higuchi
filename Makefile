GORUN=go run
GOTEST=go test -race
GOBUILD=go build -ldflags="-s -w"
BUILD_OUTPUT=bin/higuchi

serve:
	$(GORUN) . serve

test:
	$(GOTEST) ./...

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_OUTPUT)-linux-amd64 .

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_OUTPUT)-darwin-amd64 .
