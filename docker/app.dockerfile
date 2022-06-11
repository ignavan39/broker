FROM golang:1.18-alpine as builder

WORKDIR /.

ARG project

COPY go.mod ./ go.sum ./ 
RUN go mod download \
    && go mod verify

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux \
    go build -a \
    -ldflags '-extldflags "-static"' \
    -o app ./cmd/${project}/main.go

WORKDIR /

FROM scratch

COPY --from=builder /app /usr/local/bin/app
ENTRYPOINT ["/usr/local/bin/app/main"]