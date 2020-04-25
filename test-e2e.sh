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

docker run --rm -it \
    -v $HOME/.kube:/root/.kube:ro \
    -v /Users/mikulas/.minikube/profiles/minikube:/Users/mikulas/.minikube/profiles/minikube:ro \
    -v /Users/mikulas/.minikube/ca.crt:/Users/mikulas/.minikube/ca.crt:ro \
    "$IMAGE" -namespace=default -selector="app=crashloop-backoff"

#for TEST in tests/*.yaml; do
#    echo "$TEST"
#    "$KUBECTL" delete -f "$TEST"
#done
