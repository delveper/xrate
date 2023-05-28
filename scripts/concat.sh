#!/bin/bash

echo "Enter the path to the directory you want to scan:"
read path

output_file="output.ignore"

function scan_dir() {
  files=$(find "$1" -type f -name "*.go")

  # Loop through the files
  for file in $files; do
    file_name=${file##*/}

    file_path=${file%.go}

    echo "// $file_path.go" >>$output_file

    echo "$(cat $file)" >>$output_file
  done
}

scan_dir $path

cat $output_file
