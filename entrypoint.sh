#!/bin/sh
set -e

if [ ! -f "/etc/rush/conf/rush.yaml" ];then
    cp /etc/rush/rush.yaml /etc/rush/conf/rush.yaml
fi

exec "$@"
