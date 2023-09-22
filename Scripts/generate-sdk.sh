#!/bin/sh

# Define the input and output file paths
INPUT_FILE=$1
OUTPUT_DIR=$2

# Run the openapi-generator Docker image to generate the client SDK
docker run --rm \
  -v "${INPUT_FILE}:/local/openapi.json" \
  -v "${OUTPUT_DIR}:/local/out" \
  openapitools/openapi-generator-cli:v5.3.0 \
  generate \
  -i /local/openapi.json \
  -g $3 \
  -o /local/out \
  --skip-validate-spec \
  --global-property $5 \
  --additional-properties=$4

# Check if the SDK was generated successfully
if [ $? -eq 0 ]; then
  echo "SDK generated successfully in ${OUTPUT_DIR}"
else
  echo "Error generating SDK"
fi