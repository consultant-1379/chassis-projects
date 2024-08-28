#!/usr/bin/env bash
REPO_URL="selidocker.lmera.ericsson.se"
IMAGE=${REPO_URL}"/proj-nef/nef-nrfintegration"
VERSION="latest"
IMAGE_NAME=${IMAGE}:${VERSION}
docker rmi ${IMAGE_NAME}
docker build -t ${IMAGE_NAME} .
#docker push ${IMAGE_NAME}
