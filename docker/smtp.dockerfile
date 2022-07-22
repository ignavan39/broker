FROM golang:1.18-alpine as builder

WORKDIR /srv

COPY ./services/go.mod ./ ./services/go.sum ./

RUN go mod download \
    && go mod verify

COPY ./services/ ./

RUN CGO_ENABLED=0 GOOS=linux \
    go build -a \
    -ldflags '-extldflags "-static"' \
    -o smtp ./cmd/smtp/main.go

WORKDIR /srv

FROM scratch

COPY --from=builder srv/smtp /usr/local/bin/smtp
ENTRYPOINT ["/usr/local/bin/smtp/main"]