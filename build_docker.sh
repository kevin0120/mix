#!/bin/bash

version="0.1.18"

docker_repo="linshenqi/aiis"

docker build -t ${docker_repo}:${version} -t ${docker_repo}:latest .

docker push ${docker_repo}:${version}
docker push ${docker_repo}:latest
