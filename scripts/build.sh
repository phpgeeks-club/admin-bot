#!/bin/bash
export GOOS=linux
export GOARCH=amd64
go build -o geeksonator -trimpath -ldflags "-s -w"
