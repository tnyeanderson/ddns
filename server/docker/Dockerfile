FROM alpine

WORKDIR ddns

RUN apk add --no-cache dnsmasq openssh supervisor

RUN adduser -H -D ddns && passwd -d ddns

RUN mkdir -m 775 /ddns/hosts.d && chgrp ddns /ddns/hosts.d

COPY root/ /

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh /scripts/*

EXPOSE 22
EXPOSE 53

ENTRYPOINT [ "/entrypoint.sh" ]
