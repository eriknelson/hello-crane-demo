#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

$_dir/clean.sh && \
  $_dir/prep.sh && \
  $_dir/export.sh && \
  $_dir/transform.sh && \
  $_dir/apply.sh
