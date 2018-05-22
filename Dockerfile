FROM alpine:3.6

ADD ./bin/rush /rush
RUN chmod +x /rush/rush
CMD /rush/rush

