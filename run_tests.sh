#!/bin/sh

set -eu
export GO111MODULE=on
go test "${@:-./...}"
