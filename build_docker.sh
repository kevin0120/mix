#!/bin/bash

#version="0.2.56"

docker_repo="docker.pkg.github.com/masami10/rush/rush"

docker build -t ${docker_repo}:$1 -t .

docker push ${docker_repo}:$1
#docker push ${docker_repo}:latest
