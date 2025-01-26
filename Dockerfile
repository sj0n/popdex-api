ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .

FROM alpine:latest

RUN apk update && apk add ca-certificates

COPY --from=builder /run-app /usr/local/bin/
RUN chmod +x /usr/local/bin/run-app

CMD ["/usr/local/bin/run-app"]
