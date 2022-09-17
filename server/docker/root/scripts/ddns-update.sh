#!/bin/sh

set -e

ip=$1
ddns_domain="$(cat /run/ddns_domain)"

# shellcheck disable=SC1091
. /scripts/utils.sh

if [ -n "${SSH_ORIGINAL_COMMAND}" ]; then
	ip="$(echo "${SSH_ORIGINAL_COMMAND}" | awk '{print $2}')"
fi

if [ -z "${ip}" ]; then
	err "Must provide an IP address"
	exit 1
fi

if [ -z "${ddns_domain}" ]; then
	err "Must provide a DDNS_DOMAIN"
	exit 1
fi

ip_parser='^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$'

if ! echo "${ip}" | grep -q -E "${ip_parser}"; then
	err "The provided IP address is not an IP"
	exit 1
fi

echo "${ip} ${ddns_domain}" >/ddns/hosts.d/ddns
