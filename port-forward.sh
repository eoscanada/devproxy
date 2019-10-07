#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

current_dir="`pwd`"
services=(
  # First so it's easy to always have the * on it when adding/removing services
  "dgraphql-internal-v2"

  "abicodec-v2"
  "blockmeta-v2"
  "search-liverouter-v2"
  "search-archive-v2-0"
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

  let count=1
  for service in "${services[@]}"; do
    let listen_port="7000 + ${count}"
    let to_port=9000

    echo "Forwarding svc/$service (listening on $listen_port, forwarding to $to_port)"
    kubectl port-forward svc/$service $listen_port:$to_port 1> /dev/null &

    let count="$count + 1"
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
  echo "by this script on a series of ports, 7001, 7002, ..."
  echo ""
  echo "By default, these services are port forwarded:"

  for service in "${services[@]}"; do
    echo "- ${service}"
  done
}

main $@