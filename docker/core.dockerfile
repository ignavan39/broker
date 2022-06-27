FROM golang:1.18-alpine as builder

WORKDIR /.

COPY ./services/go.mod / ./services/go.sum / 
RUN go mod download \
    && go mod verify

COPY ./services/. /

RUN CGO_ENABLED=0 GOOS=linux \
    go build -a \
    -ldflags '-extldflags "-static"' \
    -o core ./cmd/core/main.go

WORKDIR /

FROM scratch

COPY --from=builder /core /usr/local/bin/core
ENTRYPOINT ["/usr/local/bin/core/main"]