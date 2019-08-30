#!/bin/bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

current_dir="`pwd`"
trap "cd \"$current_dir\"" EXIT
pushd "$ROOT" &> /dev/null

# Service definitions
SERVICES=${1:-../../service-definitions}

protoc -I$SERVICES dfuse/devproxy/v1/devproxy.proto --go_out=plugins=grpc:.
