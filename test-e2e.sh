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

function assert-selector-message() {
    SELECTOR="$1"
    MESSAGE="$2"
    docker run --rm -it \
        -v $HOME/.kube:/root/.kube:ro \
        -v /Users/mikulas/.minikube/profiles/minikube:/Users/mikulas/.minikube/profiles/minikube:ro \
        -v /Users/mikulas/.minikube/ca.crt:/Users/mikulas/.minikube/ca.crt:ro \
        "$IMAGE" -namespace=default -selector="$SELECTOR" \
    | tee .output \
    | grep --silent "$MESSAGE" \
    || {
        echo "Failed to find following message ing $SELECTOR:"
        echo "  $MESSAGE"
        echo "output:"
        cat .output
        exit 1
    }

    echo "$SELECTOR ok"
}

assert-selector-message "app=success" "Rollout successfully completed"
assert-selector-message "app=invalid-image" "x"
assert-selector-message "app=pending" "x"
assert-selector-message "app=crashloop-backoff" "Rollout failed: container main is in CrashLoopBackOff"

#for TEST in tests/*.yaml; do
#    echo "$TEST"
#    "$KUBECTL" delete -f "$TEST"
#done
