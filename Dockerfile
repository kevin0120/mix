FROM alpine:3.6

RUN apk update && apk add curl
RUN mkdir -p /etc/aiis/conf /etc/aiis/log

COPY ./entrypoint.sh /
COPY ./etc/aiis/aiis.yaml /etc/aiis/aiis.yaml
COPY ./etc/pmon/PMON.CFG /etc/aiis/PMON.CFG
COPY ./etc/pmon/ofhkht.dat /etc/aiis/ofhkht.dat
COPY ./build/aiisd /usr/bin

ENTRYPOINT ["/entrypoint.sh"]
CMD aiisd --config /etc/aiis/conf/aiis.yaml
