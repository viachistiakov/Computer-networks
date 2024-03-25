#!/bin/bash
export GOPATH="/happy/go/path/yay"
export PATH="$GOPATH/bin:$PATH"
go get github.com/mgutz/logxi/v1
imports github.com/skorobogatov/input
go install ./src/client
go install ./src/server
