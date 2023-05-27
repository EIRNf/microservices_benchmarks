#!/bin/bash

set -eEuo pipefail

# Specify the directory path
RAW_DATA_DIR=$1
START_DIR=$(pwd)
CSV_DIR="./csv_data/"
GRAPH_DIR="./graph_output/"

# Iterate over each file in the directory
for file in "$RAW_DATA_DIR"/*; do
    # Check if the current item is a file
    if [[ -f "$file" ]]; then
        # Copy raw data for processing in main dir
        cp "$file" ./
        # Get just the file name
        filename=$(basename "$file")
        bash server_raw_to_graph.sh "${filename}" "${CSV_DIR}" "${GRAPH_DIR}"
    fi
done