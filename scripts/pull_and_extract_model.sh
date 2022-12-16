#!/usr/bin/env bash

if [ $# -ne 2 ]; then
    echo "Usage: $0 <openapi_url | openapi_file> <output_file>"
    exit 1
fi

TEMP_DIR=$(mktemp -d)

echo "Downloading OpenAPI and generating go model from ${1} to ${TEMP_DIR}"
/bin/bash pull_openapi.sh $1 $TEMP_DIR
/bin/bash extract_models.sh $TEMP_DIR $2

trap 'rm -rf -- "${TEMP_DIR}" && echo "Deleted temp dir ${TEMP_DIR}" || echo "Failed to delete temp dir ${TEMP_DIR}"' EXIT 

