FROM golang:1.18-alpine as builder

WORKDIR /.

ARG project

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o app ./cmd/${project}/main.go

WORKDIR /
FROM alpine:3.14

COPY --from=builder /app /usr/local/bin/app
ENTRYPOINT ["/usr/local/bin/app/main"]