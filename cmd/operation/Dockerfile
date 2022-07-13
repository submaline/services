FROM golang:1.18-alpine as builder
WORKDIR /go/src/services

# download dependenses
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# main
COPY cmd/operation/main.go ./cmd/operation/main.go

COPY db ./db
COPY gen ./gen
COPY interceptor ./interceptor
COPY logging ./logging
COPY protocol ./protocol
COPY server ./server
COPY util ./util

RUN apk add build-base

# build as app
RUN go build -o app ./cmd/operation/main.go


FROM alpine:latest
WORKDIR /application
COPY --from=builder /go/src/services/app ./
CMD ["./app"]