#!/bin/sh

PUBSUB_PATH="$(dirname $(pwd))"

INPUT_FILE=/$PUBSUB_PATH/bcc-pubsub/src/BccCode.PubSub.Api/docs/swagger.json
OUTPUT_DIR=/$(pwd)/pubsubsdk/model

/$(pwd)/Scripts/generate-sdk.sh $INPUT_FILE $OUTPUT_DIR go "packageName=model" models
