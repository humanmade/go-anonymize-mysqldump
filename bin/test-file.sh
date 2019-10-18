#!/usr/bin/env bash

set -eu

FILE="$1"

PACKAGE_FILE=
TEST_FILE=

if [[ "$FILE" == *"_test.go"* ]]; then
  TEST_FILE="$FILE"
  PACKAGE_FILE="$(echo "$FILE" | sed -E 's/_test\.go/.go/g')"
else
  TEST_FILE="$(echo "$FILE" | sed -E 's/\.go/_test.go/g')"
  PACKAGE_FILE="$FILE"
fi

go test -v "$TEST_FILE" "$PACKAGE_FILE"
