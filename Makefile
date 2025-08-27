APP_NAME ?= ad-campaign-delivery
GO_CMD ?= go
DOCKERCMD ?= docker
GOLANGCICMD ?= golangci-lint

.PHONY: *

# needed for the performance tests
SHELL := /bin/bash

# start: runs the service in docker, remember to start docker first.
start: docker/check-engine docker/build docker/run

# generate: runs go mod tidy and generates swagger documentation
generate: tidy swagger/generate

# docker/check-engine: checks if docker engine is running
docker/check-engine:
	@if ! docker info >/dev/null 2>&1; then \
		echo "ERROR: Docker engine is not running. Please start Docker manually."; \
		exit 1; \
	fi

# docker/build: builds application image
docker/build:
	$(DOCKERCMD) build -t $(APP_NAME) .

# docker/run: starts the service container
docker/run:
	$(DOCKERCMD) run -itd --name $(APP_NAME) -v ./:/app -p 8080:8080 $(APP_NAME)

# docker/run: stops the service container
docker/stop:
	$(DOCKERCMD) stop $(APP_NAME)

# docker/run: restarts the service container
docker/restart:
	$(DOCKERCMD) restart $(APP_NAME)

# docker/test: runs unit tests in a container
docker/test:
	$(DOCKERCMD) run --rm -v ./:/app -w /app $(APP_NAME) $(GO_CMD) test -vet=off -shuffle=on -cover -covermode=count ./...

# lint: identifies lint errors
lint:
	$(DOCKERCMD) run --rm -v ./:/app -w /app $(APP_NAME) $(GOLANGCICMD) run --max-same-issues 0 --max-issues-per-linter 0

# lint/fix: fixes lint errors when possible
lint/fix:
	$(DOCKERCMD) run --rm -v ./:/app -w /app $(APP_NAME)  $(GOLANGCICMD) run --fix --max-same-issues 0 --max-issues-per-linter 0

# swagger/generate: generates swagger documentation
# run `go install github.com/swaggo/swag/cmd/swag@latest` if needed
swagger/generate:
	swag init --parseDependency --parseInternal --output docs/

# tidy: tidy dependencies
tidy:
	$(GO_CMD) mod tidy

# dev/install: development dependencies
dev/install:
	go install github.com/tsenart/vegeta/v12@latest

# performance/deliver: benchmarks performance of 'deliver endpoint'
performance/deliver:
	echo "POST http://localhost:8080/deliver" | vegeta attack \
      -duration=10s -rate=50 \
      -header "Content-Type: application/json" \
      -header "X-Consent-String: CQMGLkAQMGLkABcAKEFRBbFgAP_gAEPgAAqIJnkR_C9MQWFjcT51AfskaYxHxgACoEQgBACJgygBCAPA8IQEwGAYIAxAAqAKAAAAoiRBAAAlCAhQAAAAQAAAACCMAEAAAAAAIKBAgAARAgEACAhBGQAAEAAAAIBBABAAgAAEQBoAQBAAAAAAAAAgAAAgAACBAAAIAAAAAAEAAAAIAEgAAAAAAAAAAAAAAlAIAAAIAAAAAAAAAAAIJngAmChEQAFgQAhAAGEECABQRgAAAAAgAACBggAACAAA4AQAUGAAAAAAAAAIAAAAggABAAABAAhAAAAAQAAAAAAIAAAAAAAAACBAAAABAAAAAAgAAQAAAAAAAABAABAAgAAAABAAQBAAAAAgAAAAAAAAAACAAAAAAAAAAAEAAAAIAEAAAAAAAAAAAAAAAAAIAAAAAAAAAAAAAAAAAAA" \
      -body <(printf '{"country":"FR","device":"mobile","os":"ios"}') \
      | vegeta report

