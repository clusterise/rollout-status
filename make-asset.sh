#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

NAMESPACE="$1"
SELECTOR="$2"
OUTPUT_PATH="$3"

kubectl -n "$NAMESPACE" get deploy -o yaml -l "$SELECTOR" > "$OUTPUT_PATH/deployments.yaml"
kubectl -n "$NAMESPACE" get rs -o yaml -l "$SELECTOR" > "$OUTPUT_PATH/replicasets.yaml"
kubectl -n "$NAMESPACE" get pods -o yaml -l "$SELECTOR" > "$OUTPUT_PATH/pods.yaml"
