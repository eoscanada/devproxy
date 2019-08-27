FROM alpine:latest
RUN apk update && apk add ca-certificates curl && rm -rf /var/cache/apk/*
COPY /devproxy /app/devproxy
WORKDIR /app
