.PHONY: all dev-tools lint test latest-release publish-release
C_RED=\033[0;31m
C_GREEN=\033[0;32m
C_YELLOW=\033[0;33m
C_BLUE=\033[0;34m
NC=\033[0m

ifeq (, $(shell which go))
	@echo $(C_YELLOW)No golang in PATH, installing$(NC)
	brew install golang
endif
ifeq (, $(shell which golangci-lint))
	@echo $(C_YELLOW)No golangcli-lint in PATH, installing$(NC)
	brew install golangci/tap/golangci-lint
endif
ifeq (, $(shell which git))
	@echo $(C_YELLOW)No git in PATH, installing$(NC)
	brew install git
endif

all: lint test

lint:
	golangci-lint run

test:
	go test -mod=readonly -count=1 -p 1 -failfast -race ./...

latest-release:
	@echo "$(C_GREEN)Latest release version $(C_BLUE)$(shell git describe --tags $(shell git rev-list --tags --max-count=1))$(NC)"
	$(info)

ifneq (, $(VERSION))
ifeq (, $(shell echo $(VERSION) | grep -E "v\d+\.\d+\.\d+"))
	@echo "$(C_RED)VERSION format doesn't match expected one: vX.Y.Z (e.g. v1.2.3)$(NC)"
	@exit 1
endif
ifneq (, $(shell git status -s))
	@echo "$(C_RED)Ucommitted changes detected, cannot proceed$(NC)"
	@exit 1
endif
ifneq (main, $(shell git rev-parse --abbrev-ref HEAD))
	@git checkout main
	@git fetch --all --prune
	@git reset --hard origin/main
endif
ifneq (, $(shell git tag | grep $(VERSION)))
	@echo "$(C_RED)specified VERSION already exists$(NC)"
	@exit 1
endif

publish-release: latest-release
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
endif
