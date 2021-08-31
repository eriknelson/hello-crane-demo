#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BUILD_DIR="$_dir/build"
BIN_DIR="$_dir/bin"
PLUGIN_OUT_DIR="$BIN_DIR/plugins"
EXPORT_DIR="$_dir/export"
TRANSFORM_DIR="$_dir/transform"
OUTPUT_DIR="$_dir/output"
SRC_NAMESPACE=inventory-erik
DEST_NAMESPACE=inventory-erik-tgt
INGRESS_HOST="inventory.$DEST_NAMESPACE.10.19.2.21.nip.io"
