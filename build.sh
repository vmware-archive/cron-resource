#!/usr/bin/env bash

set -ex

go mod download

go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
ginkgo -r .

mkdir -p tmp/build

go build -o tmp/build/in in/main.go
go build -o tmp/build/check check/check.go
