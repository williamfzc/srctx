# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-alpine AS build

WORKDIR /app

# if you build it in China, add this
#ENV GOPROXY=https://goproxy.cn,direct

COPY ./go.mod ./
COPY ./go.sum ./
COPY ./Makefile ./
RUN go mod download && \
    apk add --no-cache --update make gcc g++

COPY . .
RUN make

# srctx in /app

## Deploy
FROM alpine:3

COPY --from=build /app/srctx /srctx_home/srctx
WORKDIR /srctx_home

# lsif-go
RUN wget https://github.com/sourcegraph/lsif-go/releases/download/v1.9.3/src_linux_amd64 -O lsif-go

ENTRYPOINT ["/srctx_home/srctx"]
