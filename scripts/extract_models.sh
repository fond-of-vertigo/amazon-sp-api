#! /usr/bin/env bash
# This script merges the models of the different files into one model.

if [ $# -ne 2 ]; then
    echo "Usage: $0 <go_files_dir> <output_file>"
    exit 1
fi

INPUT_FILES=$(find $1 -type f -name "model_*")
OUTPUT_FILE=$2

for file in $INPUT_FILES; do
    if [ ! -f $file ]; then
        echo "File $file does not exist."
        exit 1
    fi
    awk '/^type .* struct/{f=1} f{print; if (/}/) exit}' $file >> $OUTPUT_FILE
    echo >> $OUTPUT_FILE
done
