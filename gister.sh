#!/bin/bash

# Your GitHub Token
echo "Enter your GitHub Token"
read TOKEN

# API Endpoint
GIST_API="https://api.github.com/gists"

# Directory to scan for files
DIR=$(pwd)

# Creating a module tree
MODULE_TREE=$(tree -I 'output.md|*.git*' $DIR)

# Initializing JSON object with module_tree content
JSON_OBJECT=$(jq -n \
                  --arg mt "$MODULE_TREE" \
                  '{description: "A Gist created from files in a directory", public: true, files: {"module_tree.md": {content: $mt}}}')

# Add remaining files
while IFS= read -r -d $'\0' file; do
  FILENAME=$(basename -- "$file")
  FILE_CONTENT=$(cat "$file")
  JSON_OBJECT=$(echo "$JSON_OBJECT" | jq --arg fn "$FILENAME" --arg fc "$FILE_CONTENT" '.files[$fn] = {content: $fc}')
done < <(find $DIR \( -name '*.go' -o -name '*.md' -o -name '*.yaml' \) ! -name "README.MD" -print0)

# Creating Gist
curl -H "Authorization: token $TOKEN" -X POST -d "$JSON_OBJECT" "$GIST_API"

# Creating Gist and getting response
RESPONSE=$(curl -H "Authorization: token $TOKEN" -X POST -d "$JSON_OBJECT" "$GIST_API")
# Extracting and printing raw URLs
echo "$RESPONSE" | jq -r '.files[] | .raw_url' >> output_links.ignore