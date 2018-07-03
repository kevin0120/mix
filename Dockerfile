FROM alpine:3.6

RUN apk update && apk add curl
RUN mkdir -p /etc/rush/conf /etc/rush/log

COPY ./entrypoint.sh /
COPY ./etc/rush/rush.yaml /etc/rush/rush.yaml
COPY ./etc/rush/api.json /etc/rush
COPY ./build/rushd /usr/bin

ENTRYPOINT ["/entrypoint.sh"]
CMD rushd --config /etc/rush/conf/rush.yaml
