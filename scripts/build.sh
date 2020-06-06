#!/bin/bash
ROOT_DIR=$PWD

# build all commands.
for file in `find "${PWD}"/golang -name "main.go" -type f`; do
    dir=`dirname "${file}"`
    basename=`basename "${dir}"`
    cd "${dir}"
    go build -o "${ROOT_DIR}"/cmd/"${basename}"
done
