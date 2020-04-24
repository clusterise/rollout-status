#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

IMAGE="$1"

#minikube start --kubernetes-version=v1.16.9 --driver=virtualbox
#KUBECTL=kubectl-1.16.3
#
#for TEST in tests/*.yaml; do
#    echo "$TEST"
#    "$KUBECTL" apply -f "$TEST"
#done

docker run --rm -it "$IMAGE" -namespace=default -selector="app=crashloop-backoff"

#for TEST in tests/*.yaml; do
#    echo "$TEST"
#    "$KUBECTL" delete -f "$TEST"
#done
