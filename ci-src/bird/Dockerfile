FROM alpine:latest
RUN apk add --no-cache bird
EXPOSE 179/tcp
ENTRYPOINT ["bird", "-d"]