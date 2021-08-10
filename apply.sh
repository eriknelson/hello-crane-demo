#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source $_dir/var.sh

./bin/crane apply \
  --export-dir=$_dir/export/resources \
  --transform-dir=$_dir/transform \
  --output-dir=$_dir/output
