#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
kubectl create -f $_dir/argo-k8s/00-Namespace*
./bin/crane transfer-pvc \
  --skip-quiesce --local \
  --source-context=ie \
  --source-pvc-name=data-inventory-postgresql-0 \
  --destination-pvc-namespace=inventory-erik-tgt
