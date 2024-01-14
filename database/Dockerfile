FROM alpine:latest

RUN apk add --no-cache sqlite

RUN mkdir /database

WORKDIR /database

CMD ["sh", "-c", "while true ; do sleep 10 ; done"]