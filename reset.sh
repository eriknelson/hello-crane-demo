#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

oc delete -f $_dir/argo
$_dir/clean.sh
