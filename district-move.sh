#!/bin/bash

# Check if required arguments are provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <max_number> <directory_prefix>"
    echo "Example: $0 8 NC"
    exit 1
fi

max_number=$1
dir_prefix=$2

# Loop through numbers 1 to max_number
for i in $(seq 1 "$max_number"); do
    # Define source and destination paths
    source_dir="${dir_prefix}-${i}/shape.geojson"
    dest_file="${i}.geojson"
    
    # Check if source file exists
    if [ -f "$source_dir" ]; then
        # Move and rename the file
        mv "$source_dir" "$dest_file"
        echo "Moved ${dir_prefix}-${i}/shape.geojson to ${i}.geojson"
    else
        echo "Warning: ${source_dir} not found"
    fi
done