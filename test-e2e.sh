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
    assert-selector-message-ns "$1" "default" "$2"
}

function assert-selector-message-ns() {
    SELECTOR="$1"
    NAMESPACE="$2"
    MESSAGE="$3"
    { docker run --rm -it \
        -v $HOME/.kube:/root/.kube:ro \
        -v /Users/mikulas/.minikube/profiles/minikube:/Users/mikulas/.minikube/profiles/minikube:ro \
        -v /Users/mikulas/.minikube/ca.crt:/Users/mikulas/.minikube/ca.crt:ro \
        "$IMAGE" -namespace="$NAMESPACE" -selector="$SELECTOR" || true
    } \
    | tee .output \
    | grep --silent --fixed-strings "$MESSAGE" \
    || {
        echo "Failed to find following message in $SELECTOR:"
        echo "  $MESSAGE"
        echo "output:"
        cat .output
        exit 1
    }

    echo "$SELECTOR ok"
}

assert-selector-message "app=not-found" \
    'Selector "app=not-found" did not match any Deployments'
assert-selector-message "app=success" \
    'Rollout successfully completed'
assert-selector-message-ns "app=limit-range" "limit-range" \
    'Rollout failed: replicaset limit-range-7dfcd777fd failed to create pods: Pod "limit-range-7dfcd777fd-f99jf" is invalid: spec.containers[0].resources.requests: Invalid value: "200Mi": must be less than or equal to memory limit'
assert-selector-message-ns "app=resource-quota" "resource-quota" \
    'Rollout failed: replicaset resource-quota-6884c5558d failed to create pods: pods "resource-quota-6884c5558d-bc5pq" is forbidden: exceeded quota: main, requested: memory=200Mi, used: memory=0, limited: memory=100Mi'
assert-selector-message "app=invalid-image" \
    'Rollout failed: container main is in ImagePullBackOff: Back-off pulling image "bogus-image:does-not-exist"'
assert-selector-message "app=pending" \
    'Rollout failed: failed to scheduled pod: 0/1 nodes are available: 1 Insufficient memory.'
assert-selector-message "app=crashloop-backoff" \
    'Rollout failed: container main is in CrashLoopBackOff'
assert-selector-message "app=readiness" \
    'Rollout failed: deployment limit-range is not progressing: ReplicaSet "limit-range-7dfcd777fd" has timed out progressing.'

#for TEST in tests/*.yaml; do
#    echo "$TEST"
#    "$KUBECTL" delete -f "$TEST"
#done
