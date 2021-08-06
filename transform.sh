#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

./bin/crane transform \
  --export-dir=$_dir/export/resources \
  --plugin-dir=$_dir/bin/plugins \
  --transform-dir=$_dir/transform \
  --optional-flags="dest-namespace=$DEST_NAMESPACE"
