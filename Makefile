GOTEST=go test
GOBUILD=go build
BINARY_OUTPUT=bin/higuchi

test:
	$(GOTEST) ./...

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_OUTPUT)-linux-amd64 .

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_OUTPUT)-darwin-amd64 .
