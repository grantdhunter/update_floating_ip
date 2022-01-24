FROM golang:1.17-alpine3.15 as builder

WORKDIR /build

COPY go.mod /build
COPY go.sum /build
RUN go mod download

COPY *.go /build

RUN go build

FROM alpine:3.15
WORKDIR /
COPY --from=builder /build/update_floating_ip .




ENTRYPOINT ["/update_floating_ip"]