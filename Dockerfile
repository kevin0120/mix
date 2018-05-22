FROM alpine:3.6

copy ./bin/rush /rush/
cmd rush

