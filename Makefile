all:

dev-tools:
ifeq (, $(shell which go))
	$(warning "No golang in PATH, consider doing brew install golang")
	brew install golang
endif
ifeq (, $(shell which golangci-lint))
	brew install golangci/tap/golangci-lint
endif

lint:
	golangci-lint run

test:
	go test -mod=readonly -count=1 -p 1 -failfast -race ./...
