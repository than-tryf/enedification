#!/usr/bin/env bash


#echo $PWD
docker run --rm -v "$GOPATH/src/":/go/src/ -w /go/src/enedification/datacenter-agent golang:1.11 go build -v
#docker run --rm -v "$GOPATH/src/":/go/src/ -w /go/src/github.com/enedification/datacenter-agent golang:1.11 go build -v
