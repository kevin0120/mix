#!/bin/bash

version="0.2.20"

docker_repo="linshenqi/rush"

docker build -t ${docker_repo}:${version} -t ${docker_repo}:latest .

docker push ${docker_repo}:${version}
docker push ${docker_repo}:latest
