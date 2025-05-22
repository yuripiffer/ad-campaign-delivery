APP_NAME ?= ad-campaign-delivery
GO_CMD ?= go
DOCKERCMD ?= docker
GOLANGCICMD ?= golangci-lint

.PHONY: *

start: docker/check-engine docker/build docker/run

generate: tidy swagger/generate

docker/check-engine:
	@if ! docker info >/dev/null 2>&1; then \
		echo "ERROR: Docker engine is not running. Please start Docker manually."; \
		exit 1; \
	fi

docker/build:
	$(DOCKERCMD) build -t $(APP_NAME) .

docker/run:
	$(DOCKERCMD) run -itd --name $(APP_NAME) -v ./:/app -p 8080:8080 $(APP_NAME)

docker/test:
	$(DOCKERCMD) run --rm -v ./:/app -w /app $(APP_NAME) $(GO_CMD) test -vet=off -shuffle=on -cover -covermode=count ./...

lint:
	$(DOCKERCMD) run --rm -v ./:/app -w /app $(APP_NAME) $(GOLANGCICMD) run --max-same-issues 0 --max-issues-per-linter 0

lint/fix:
	$(DOCKERCMD) run --rm -v ./:/app -w /app $(APP_NAME)  $(GOLANGCICMD) run --fix --max-same-issues 0 --max-issues-per-linter 0

#run `go install github.com/swaggo/swag/cmd/swag@latest` if needed
swagger/generate:
	swag init --parseDependency --parseInternal --output docs/

tidy:
	$(GO_CMD) mod tidy