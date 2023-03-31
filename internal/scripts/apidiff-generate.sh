#!/bin/bash
# Usage: apidiff-generate.sh <module folder> <destination folder>
# 
# Generates apidiff-compatible gcexport data files 
# for all non-internal packages in the module.

set -eo pipefail

BASE_DIR=$(pwd)
pushd "$1" > /dev/null
for pkg in $(go list ./...); do
    if [[ "$pkg" =~ .*"/internal/".* ]]; then
        # Internal packages don't have breaking changes
        continue
    fi

    # Write gcexport data into $2/<package path>/apidiff.state
    OUT_DIR=$BASE_DIR/$2/$pkg
    mkdir -p "$OUT_DIR"
    apidiff -w "$OUT_DIR/apidiff.state" "$pkg"
done
popd > /dev/null
