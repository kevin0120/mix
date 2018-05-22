FROM alpine:3.6

copy ./bin/rush /rush/
run chmod +x /rush/rush
cmd /rush/rush

