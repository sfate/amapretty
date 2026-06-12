.PHONY: all dev-tools lint audit test latest-release publish-release
C_RED=\033[0;31m
C_GREEN=\033[0;32m
C_YELLOW=\033[0;33m
C_BLUE=\033[0;34m
NC=\033[0m

all: lint audit test

dev-tools:
	@command -v go >/dev/null || (echo "$(C_RED)go is required but was not found in PATH$(NC)" && exit 1)
	@command -v git >/dev/null || (echo "$(C_RED)git is required but was not found in PATH$(NC)" && exit 1)

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0 run --allow-parallel-runners

audit:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

test:
	go test -mod=readonly -count=1 -p 1 -failfast -race ./...

latest-release:
	@echo "$(C_GREEN)Latest release version $(C_BLUE)$(shell git describe --tags $(shell git rev-list --tags --max-count=1))$(NC)"
	$(info)

publish-release: latest-release
	@test -n "$(VERSION)" || (echo "$(C_RED)VERSION is required, e.g. make publish-release VERSION=v1.2.3$(NC)" && exit 1)
	@echo "$(VERSION)" | grep -Eq "^v[0-9]+\.[0-9]+\.[0-9]+$$" || (echo "$(C_RED)VERSION format doesn't match expected one: vX.Y.Z (e.g. v1.2.3)$(NC)" && exit 1)
	@test -z "$$(git status -s)" || (echo "$(C_RED)Uncommitted changes detected, cannot proceed$(NC)" && exit 1)
	@test "$$(git rev-parse --abbrev-ref HEAD)" = "main" || (echo "$(C_RED)publish-release must be run from main$(NC)" && exit 1)
	@! git tag | grep -Fx "$(VERSION)" >/dev/null || (echo "$(C_RED)specified VERSION already exists$(NC)" && exit 1)
	@echo "$(C_YELLOW)Clean up$(NC)"
	@go mod tidy
	@echo "$(C_YELLOW)Run lint$(NC)"
	@$(MAKE) lint
	@echo "$(C_YELLOW)Run tests$(NC)"
	@$(MAKE) test
	@echo "$(C_YELLOW)Tag release$(NC)"
	@git tag $(VERSION)
	@echo "$(C_YELLOW)Push release$(NC)"
	@git push origin HEAD > /dev/null
	@echo "$(C_YELLOW)Push release tag$(NC)"
	@git push origin $(VERSION) > /dev/null
	@echo "$(C_GREEN)New version $(VERSION) released$(NC)"
