#!/bin/bash

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: $0 <model-name> <path-to-adapter-gguf>"
    exit 1
fi

model_name="$1"
adapter_path="$2"

if [ ! -f "$adapter_path" ]; then
    echo "Error: Adapter not found at the specified path: $adapter_path"
    exit 1
fi

# Get the IP address of the ollama container
container_ip=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ollama)

# Check if container IP retrieval was successful
if [ -z "$container_ip" ]; then
    echo "Error: Unable to retrieve IP address of the ollama container."
    exit 1
fi

BASE_URL="http://$container_ip:11434/api"

# Pull the model
status_code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/pull" -d '{"model": "llama3:8b-instruct-q4_0"}')
if [ "$status_code" -eq 200 ]; then
  echo "Downloaded llama model"
else
  echo "Request failed with status code: $status_code"
fi

# Get adapter checksum
digest=$(sha256sum "$adapter_path" | cut -d' ' -f1)
if [ -z "$digest" ]; then
    echo "Error: Unable to retrieve digest of adapter"
    exit 1
fi

# Upload adapter as blob
status_code=$(curl -X POST -T "$adapter_path" -s -o /dev/null -w "%{http_code}" "$BASE_URL/blobs/sha256:$digest")

if [ "$status_code" -eq 201 ] || [ "$status_code" -eq 200 ]; then
  echo "Uploaded adapter to ollama"
else
  echo "Adapter upload request failed with status code: $status_code"
  exit 1
fi

# Read the contents of Modelfile
#modelfile_content=$(sed 's/\\/\\\\/g; s/"/\\"/g; s/$/\\n/' $"modelfile_path"| tr -d '\n')
modelfile_content="FROM llama3:8b-instruct-q4_0\\nADAPTER /root/.ollama/models/blobs/sha256-$digest"

status_code=$(curl -s -X POST -s -o /dev/null -w "%{http_code}" "$BASE_URL/create" -d '{"name": "'"$model_name"'", "modelfile": "'"$modelfile_content"'"}')
if [ "$status_code" -eq 201 ] || [ "$status_code" -eq 200 ]; then
  echo "Model creation successful"
else
  echo "Model creation request failed with status code: $status_code"
  exit 1
fi
