ARG GO_VERSION=1
FROM alpine:latest AS builder

RUN apk add --no-cache --update go gcc g++
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=1 go build -o /echo-api .

FROM alpine:latest

RUN apk update && apk add ca-certificates
COPY --from=builder /echo-api /usr/local/bin/echo-api

CMD ["echo-api"]