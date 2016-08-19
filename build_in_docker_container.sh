#!/usr/bin/bash
set -ex

cd $GOPATH/src/github.com/pivotal-cf-experimental/cron-resource
go get github.com/tools/godep
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
godep restore
ginkgo -r .
go build -o built-in in/main.go
go build -o built-check check/check.go
