#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source $_dir/var.sh

$_dir/clean.sh
kubectl delete -f $_dir/argo-k8s
kubectl delete pod -n $SRC_NAMESPACE --selector=app=crane2
