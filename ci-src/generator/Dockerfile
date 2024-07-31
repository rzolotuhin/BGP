FROM alpine:latest

ARG TARGETOS
ARG TARGETARCH

RUN apk add libc6-compat

ENV listPath=/run/bird/ru-subnet.list
WORKDIR /usr/bin
COPY bin/generator_${TARGETARCH}_${TARGETOS} generator
ENTRYPOINT ["generator"]
