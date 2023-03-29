all:

dev-tools:
ifeq (, $(shell which go))
	$(warning "No golang in PATH, consider doing brew install golang")
	brew install golang
endif
ifeq (, $(shell which golangci-lint))
	brew install golangci/tap/golangci-lint
endif
ifeq (, $(shell which git))
	brew install git
endif

lint: dev-tools
	golangci-lint run

test: dev-tools
	go test -mod=readonly -count=1 -p 1 -failfast -race ./...

check-release: dev-tools
ifneq (, $(shell git status -s))
	$(error Ucommitted changes detected, cannot proceed)
endif
ifeq (, $(shell echo $(VERSION) | grep -E "v\d+\.\d+\.\d+"))
	$(error VERSION format doesn't match expected one: vX.Y.Z (e.g. v1.2.3))
endif
ifeq ($(shell git describe --tags $(shell git rev-list --tags --max-count=1)), $(shell echo $(VERSION) | grep -E "v\d+\.\d+\.\d+"))
	$(error specified VERSION already exists)
endif
$(info Latest release version $(shell git describe --tags $(shell git rev-list --tags --max-count=1)))

publish-release: check-release
	$(info Clean up)
	@go mod tidy
	$(info Run tests)
	@go test -mod=readonly -count=1 -p 1 -failfast -race ./...
	$(info Tag release)
	@git tag $(VERSION)
	$(info Push release)
	@git push origin HEAD > /dev/null
	$(info Push release tag)
	@git push origin $(VERSION) > /dev/null
	$(info New version $(VERSION) released)

latest-release: dev-tools
	@git describe --tags $(shell git rev-list --tags --max-count=1)
