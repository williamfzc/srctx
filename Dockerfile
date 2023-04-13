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
WORKDIR /srctx_home
# Install git
RUN apk add --no-cache git bash
# lsif-go
RUN wget https://github.com/sourcegraph/lsif-go/releases/download/v1.9.3/src_linux_amd64 -O lsif-go
# golang
COPY --from=golang:1.19-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

COPY --from=build /app/srctx /srctx_home/srctx
COPY ./scripts/create_and_diff.sh /srctx_home/create_and_diff.sh

RUN chmod +x /srctx_home/create_and_diff.sh && \
    chmod +x /srctx_home/lsif-go && \
    git config --global --add safe.directory /app

ENV PATH="${PATH}:/srctx_home"

ENTRYPOINT ["/srctx_home/create_and_diff.sh"]
