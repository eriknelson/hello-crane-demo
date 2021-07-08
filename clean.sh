#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source $_dir/var.sh
rm -rf $BUILD_DIR
rm -rf $BIN_DIR
