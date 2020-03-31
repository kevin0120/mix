#!/bin/sh
set -e

if [ ! -f "/etc/aiis/conf/aiis.yaml" ];then
    cp /etc/aiis/aiis.yaml /etc/aiis/conf/aiis.yaml
fi

if [ ! -f "/etc/aiis/conf/PMON.CFG" ];then
    cp /etc/aiis/PMON.CFG /etc/aiis/conf/PMON.CFG
fi

if [ ! -f "/etc/aiis/conf/ofhkht.dat" ];then
    cp /etc/aiis/ofhkht.dat /etc/aiis/conf/ofhkht.dat
fi

exec "$@"
