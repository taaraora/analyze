#!/bin/bash
# Exit script when command fails
set -o errexit
# if any of the commands in pipeline fails, script will exit
set -o pipefail

# log into docker
echo "Log into docker"
docker login -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD"
# push to docker
echo "Pushing to Docker $DOCKER_IMAGE_NAME:$DOCKER_IMAGE_TAG"
docker push "$DOCKER_IMAGE_NAME":"$DOCKER_IMAGE_TAG"
echo "Pushing to Docker $NODEAGENT_DOCKER_IMAGE_NAME:$DOCKER_IMAGE_TAG"
docker push "$NODEAGENT_DOCKER_IMAGE_NAME":"$DOCKER_IMAGE_TAG"
echo "Pushing to Docker $JOB_DOCKER_IMAGE_NAME:$DOCKER_IMAGE_TAG"
docker push "$JOB_DOCKER_IMAGE_NAME":"$DOCKER_IMAGE_TAG"

if [[ "$TRAVIS_TAG" =~ ^v[0-9]. ]]; then
    docker push "$DOCKER_IMAGE_NAME":"latest"
    docker push "$NODEAGENT_DOCKER_IMAGE_NAME":"latest"
    docker push "$JOB_DOCKER_IMAGE_NAME":"latest"
fi
