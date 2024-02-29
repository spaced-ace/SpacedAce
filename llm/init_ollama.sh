#!/bin/bash

# Get the IP address of the ollama container
container_ip=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ollama)

# Check if container IP retrieval was successful
if [ -z "$container_ip" ]; then
    echo "Error: Unable to retrieve IP address of the ollama container."
    exit 1
fi

# Generate the URL for the curl request
url="http://$container_ip:11434/api/create"

# Check if Modelfile exists in the current directory
if [ ! -f "Modelfile" ]; then
    echo "Error: Modelfile not found in the current directory."
    exit 1
fi

# Read the contents of Modelfile
modelfile_content=$(sed 's/\\/\\\\/g; s/"/\\"/g; s/$/\\n/' Modelfile | tr -d '\n')

curl -s -X POST "$url" -d '{"name": "mistral-quizgen", "modelfile": "'"$modelfile_content"'"}'
