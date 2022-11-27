#!/bin/bash 

set -x  # optional

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
devbase=$SCRIPT_DIR
port=6060

docker run \
    --rm \
    -e "GOPATH=/tmp/go" \
    -p 127.0.0.1:$port:$port \
    -v $devbase:/tmp/go/src/ \
    --name godoc \
    golang \
    bash -c "go install golang.org/x/tools/cmd/godoc@latest ; echo http://localhost:$port/pkg/ ; cd /tmp/go/src ; /tmp/go/bin/godoc -http=:$port"