FROM alpine

WORKDIR /ddns

RUN apk add --update --no-cache openssh curl

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]

