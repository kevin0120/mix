#!/bin/bash

while read -r line; do declare  $line; done < ./version.env
_tag=`echo $line|awk -F '=' '{print $2}'`
echo $_tag;

docker_repo="registry.cn-hangzhou.aliyuncs.com"

docker build -t ${docker_repo}:rush_${_tag} -t ${docker_repo}:rush_latest .

docker push ${docker_repo}:rush_${_tag}
docker push ${docker_repo}:rush_latest
