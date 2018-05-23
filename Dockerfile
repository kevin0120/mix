FROM alpine:3.6

ADD ./bin/rush /rush/rush
CMD /rush/rush

