#!/bin/sh

set -e

# shellcheck disable=SC1091
. /scripts/utils.sh

authorized_keys=/etc/ssh/authorized_keys

if [ -z "${DDNS_DOMAIN}" ]; then
	err "The domain being controlled by DDNS must be set with DDNS_DOMAIN"
	exit 1
fi

# Have to save this to a file so it can be accessed by the ssh user
echo "${DDNS_DOMAIN}" >/run/ddns_domain

### Set up ssh

# Generate host keys
ssh-keygen -A

# Only use the first line of the public key file for safety
public_key="$(head -n 1 /ddns/ssh.key.pub)"

mkdir -p "${authorized_keys}"

cat >"${authorized_keys}/ddns" <<-EOF
	command="/scripts/ddns-update.sh" ${public_key}
EOF

### Start services
supervisord \
	--configuration /etc/supervisor/supervisord.conf \
	--pidfile /run/supervisord.pid
