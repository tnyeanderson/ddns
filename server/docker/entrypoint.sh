#!/bin/sh

set -e

# shellcheck disable=SC1091
. /scripts/utils.sh

authorized_keys=/etc/ssh/authorized_keys
hosts_dir=/ddns/hosts.d
ssh_host_key=/ddns/ssh-host-keys/ssh_host_ed25519_key

if [ -z "${DDNS_DOMAIN}" ]; then
	err "The domain being controlled by DDNS must be set with DDNS_DOMAIN"
	exit 1
fi

# Have to save this to a file so it can be accessed by the ssh user
echo "${DDNS_DOMAIN}" >/run/ddns_domain

### Set up ssh

# Generate host key
if [ ! -f "${ssh_host_key}" ]; then
	mkdir -p "$(dirname "${ssh_host_key}")"
	ssh-keygen -t ed25519 -C ddns-host-key -f "${ssh_host_key}" -N ''
fi

# Only use the first line of the public key file for safety
public_key="$(head -n 1 /ddns/ssh.key.pub)"

mkdir -p "${authorized_keys}"

# Set up hosts.d
mkdir -p "${hosts_dir}"
chmod 775 "${hosts_dir}"
chgrp ddns "${hosts_dir}"

cat >"${authorized_keys}/ddns" <<-EOF
	command="/scripts/ddns-update.sh" ${public_key}
EOF

### Start services
supervisord \
	--configuration /etc/supervisor/supervisord.conf \
	--pidfile /run/supervisord.pid
