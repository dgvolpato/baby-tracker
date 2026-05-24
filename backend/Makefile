BINARY = cmd/api/bootstrap

.PHONY: build deploy local clean

build:
	GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -o $(BINARY) ./cmd/api

deploy: build
	sam deploy --guided

local:
	go run ./cmd/local

clean:
	rm -f $(BINARY)
