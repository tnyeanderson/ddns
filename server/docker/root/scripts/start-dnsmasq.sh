#!/bin/sh

set -e

# shellcheck disable=SC1091
. /scripts/utils.sh

if [ -z "${DDNS_DOMAIN}" ]; then
	err "The domain being controlled by DDNS must be set with DDNS_DOMAIN"
	exit 1
fi

dnsmasq \
	--keep-in-foreground \
	--no-resolv \
	--no-hosts \
	--hostsdir=/ddns/hosts.d \
	--cache-size=0 \
	--max-ttl=60 \
	--auth-zone "${DDNS_DOMAIN}" \
	--auth-server "${DDNS_DOMAIN}" \
	--log-facility=- \
	--log-queries=extra
