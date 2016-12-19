#!/usr/bin/env bash
#
# This script builds the application from source.
set -e

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that directory
cd $DIR

if [ "$(go env GOOS)" = "freebsd" ]; then
  export CC="clang"
fi

# Build
echo "--> Building..."
go build \
    -o bin/datalock \
    -v
