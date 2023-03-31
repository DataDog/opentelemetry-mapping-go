#!/bin/bash
# Usage: apidiff-compare.sh <module folder> <gcexport data folder>
# 
# Compares current API with saved API in gcexport data files
# for all non-internal packages in the module.
set -eo pipefail

BASE_DIR=$(pwd)
pushd "$1" > /dev/null
for pkg in $(go list ./...); do
    if [[ "$pkg" =~ .*"/internal/".* ]]; then
        # Internal packages don't have breaking changes
        continue
    fi
    GCEXPORT_DIR=$BASE_DIR/$2/$pkg
    changes=$(apidiff "$GCEXPORT_DIR/apidiff.state" "$pkg")
    if [[ -n "$changes" ]]; then
        echo "Changes for $pkg:"
        echo "$changes"
    fi

    # Search for "Incompatible changes:" header
    if [[ "$changes" =~ .*"Incompatible changes:".* ]]; then
        echo "Incompatible changes found!"
        exit 1
    fi
done
popd > /dev/null
