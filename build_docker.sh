#!/bin/bash

while read -r line; do declare  $line; done < ./version.env
_tag=`echo $line|awk -F '=' '{print $2}'`
echo $_tag;

docker build -t registry.us-west-1.aliyuncs.com/us/docker:rush_${_tag} -t registry.us-west-1.aliyuncs.com/us/docker:rush_latest .

docker push registry.us-west-1.aliyuncs.com/us/docker:rush_${_tag}
docker push registry.us-west-1.aliyuncs.com/us/docker:rush_latest
