#!/bin/bash

version="0.0.1"

docker_repo="registry.cn-hangzhou.aliyuncs.com/ca/docker"

docker build -t ${docker_repo}:rush_${version} -t ${docker_repo}:rush_latest .

docker push ${docker_repo}:rush_${version}
docker push ${docker_repo}:rush_latest
