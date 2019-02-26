#!/bin/bash
# Exit script when command fails
set -o errexit
# Exit script when it tries to use undeclared variables
set -o nounset
# if any of the commands in pipeline fails, script will exit
set -o pipefail

echo "$TRAVIS_REPO_SLUG":"$TAG"
# build the docker container
echo "Building Docker container"
make gen-assets && make build-image

if [ $? -eq 0 ]; then
	echo "Complete"
else
	echo "Build Failed"
	exit 1
fi
