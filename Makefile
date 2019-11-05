include Makefile.ledger
all: lint install docker

install: go.sum
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/famd
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/famcli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

docker:
	docker build -t fam .

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

test:
	@go test -mod=readonly $(PACKAGES)
