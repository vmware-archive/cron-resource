#!/usr/bin/env bash

set -ex

go mod download

go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
ginkgo -r .

go build -o built-in in/main.go
go build -o built-check check/check.go
