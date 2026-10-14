#!/bin/sh
# Copies the delivery token test vectors from a bcc-pubsub checkout next to this repository
# (../bcc-pubsub/src/sdk/test-vectors/delivery-token). Run after the vectors change there.

set -e
SOURCE="$(dirname "$(pwd)")/bcc-pubsub/src/sdk/test-vectors/delivery-token"
TARGET="$(pwd)/pubsubsdk/testdata/delivery-token"

rm -f "$TARGET"/*.json
cp "$SOURCE"/*.json "$TARGET"/
echo "Copied $(ls "$TARGET"/*.json | wc -l) vectors from $SOURCE"
