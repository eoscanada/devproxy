#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

current_dir="`pwd`"
services=(
  # The internal service requires a `*` (see `main.go` default flags), put it first so it's easy when adding new service
  "dgraphql-internal-v2:7001:9000"

  # gRPC services
  "abicodec-v2:7002:9000"
  "blockmeta-v2:7003:9000"
  "search-liverouter-v2:7004:9000"
  "search-archive-v2-0:7005:9000"

  # HTTP proxy services
  "fluxdb-server-v2:8080:80"
  "nodeos-api-v2:9999"
)

teardown() {
  for job in `jobs -p`; do
    kill -s TERM $job &> /dev/null || true
  done
}

main() {
  if [[ $1 == "--help" || $1 == "-h" ]]; then
    usage
    exit
  fi

  trap teardown EXIT
  pushd "$ROOT" &> /dev/null

  if [[ $1 != "" ]]; then
    services=$1; shift
  fi

  for service in "${services[@]}"; do
    name=$(printf $service | cut -f1 -d':')
    listen_port=$(printf $service | cut -f2 -d':')
    to_port=$(printf $service | cut -f3 -d':')
    if [[ $to_port == "" ]]; then
      to_port=$listen_port
    fi

    echo "Forwarding svc/$name (listening on $listen_port, forwarding to $to_port)"
    kubectl port-forward svc/$name $listen_port:$to_port 1> /dev/null &
  done

  echo ""
  echo "Press Ctrl+C to terminal all port forwarding"
  for job in `jobs -p`; do
    wait $job || true
  done
}

usage() {
  echo "usage: port-forward.sh [<services>]"
  echo ""
  echo "For development purposes, start port-forwarding of all services known"
  echo "to the port specified by the service."
  echo ""
  echo "By default, these services are port forwarded:"

  for service in "${services[@]}"; do
    echo "- ${service}"
  done
}

main $@