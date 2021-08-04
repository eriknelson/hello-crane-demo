#!/bin/bash
_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
./bin/crane export --namespace=nginx-example --export-dir=export
