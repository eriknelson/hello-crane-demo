#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source $_dir/var.sh

mkdir -p $BUILD_DIR
mkdir -p $BIN_DIR
mkdir -p $PLUGIN_OUT_DIR

git clone https://github.com/konveyor/crane.git $BUILD_DIR/crane
pushd $BUILD_DIR/crane
go build -o $BIN_DIR/crane
popd

plugins=$(ls -1 $_dir/plugins)
shopt -s nullglob
plugins=($_dir/plugins/*/)
shopt -u nullglob
for pluginDir in "${plugins[@]}"; do
  pluginName=$(echo $p | perl -ne '/.*\/(.*)\/$/ && print "$1"')
  pushd $pluginDir
  go build -o $PLUGIN_OUT_DIR/$pluginName $pluginDir
  popd
done
