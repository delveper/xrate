#!/bin/bash

env_file=$1

if [ ! -f $env_file ]; then
  echo "Error: .env file does not exist"
  exit 1
fi

while read -r var val; do
  export $var=$val
done < $env_file
