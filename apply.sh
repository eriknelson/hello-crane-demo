#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

oc new-project $DEST_NAMESPACE

./bin/crane apply \
  --export-dir=$_dir/export/resources \
  --transform-dir=$_dir/transform \
  --output-dir=$_dir/output
