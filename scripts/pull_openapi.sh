#! /usr/bin/env bash

if [[ -z "$1" ]]; then
  # argument is null or empty
  echo "No input provided"
  exit 1
fi
if [[ -z "$2" ]]; then
  # argument is null or empty
  echo "No output provided, with respect to current working directory"
  exit 1
fi

docker run --rm \
  -v $2:/local openapitools/openapi-generator-cli generate \
  -i $1 \
  -g go \
  -o /local/
