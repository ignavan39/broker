FROM golang:1.18-alpine as builder

WORKDIR /srv/app

RUN apk update \
 && apk upgrade \
 && apk add --no-cache npm \
 && npm install --global nodemon

ARG project
ENV PROJECT=./cmd/${project}/main.go
COPY . ./

RUN go mod download

ENTRYPOINT ["nodemon", "--watch", "./**/*.go", "--signal", "SIGTERM", "--exec", "go", "run", "$PROJECT"]
