#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source $_dir/var.sh

mkdir -p $BUILD_DIR
mkdir -p $BIN_DIR
mkdir -p $PLUGIN_OUT_DIR

#git clone --depth=1 https://github.com/konveyor/crane.git $BUILD_DIR/crane
#cd $BUILD_DIR/crane
#go build -o $BIN_DIR/crane

plugins=$(ls -1 $_dir/plugins)
shopt -s nullglob
plugins=($_dir/plugins/*/)
shopt -u nullglob
for p in "${plugins[@]}"; do
  pluginMain="$(ls $p)"
  cd $p
  go build -o $PLUGIN_OUT_DIR $pluginMain
  popd
done
