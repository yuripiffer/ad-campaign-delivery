FROM golang:1.24.2-alpine

RUN apk --no-cache update && apk upgrade && \
    apk add --no-cache musl-dev curl alpine-sdk linux-headers ca-certificates bash dumb-init tzdata git make gcc

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["air"]
