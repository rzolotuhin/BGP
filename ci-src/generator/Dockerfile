FROM alpine:latest

RUN apk add libc6-compat

ENV listPath /run/bird/ru-subnet.list
WORKDIR /usr/bin
COPY bin/generator_amd64_linux generator
ENTRYPOINT ["generator"]
