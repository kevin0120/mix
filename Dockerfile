FROM alpine:3.6

RUN mkdir -p /etc/aiis/conf /etc/aiis/log

COPY ./entrypoint.sh /
COPY ./etc/aiis/aiis.yaml /etc/aiis/aiis.yaml
COPY ./etc/pmon/PMON.CFG /etc/aiis/PMON.CFG
COPY ./etc/aiis/api.json /etc/aiis
COPY ./build/aiisd /usr/bin

ENTRYPOINT ["/entrypoint.sh"]
CMD aiisd --config /etc/aiis/conf/aiis.yaml
