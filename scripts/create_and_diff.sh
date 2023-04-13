#!/bin/sh
set -e

cd /app
lsif-go -v
srctx diff
