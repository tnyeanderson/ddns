#!/bin/sh

set -e

### CONFIGURATION: Set these as environment variables
#SSH_HOST="bar.com"
#SSH_PORT="2222"
#
# Optional:
#FORCE=yes

new_ip=$1

# Set value to "yes" to make update the DDNS entry
# even if it hasn't changed
# shellcheck disable=SC2153
force=$FORCE

if [ -n "${new_ip}" ]; then
	# Always force if an IP was provided directly
	echo "IP was provided directly: ${new_ip}"
	force=yes
fi

remote_user=ddns
cache_file=/ddns/current-ip
ssh_keyfile=${SSH_PRIVATE_KEY:-/ddns/ssh.key}
ssh_known_hosts=/ddns/known_hosts/ssh_host_ed25519_key.pub

ssh_port=${SSH_PORT:-22}

debug() {
	if [ "${DEBUG}" = 'true' ]; then
		echo "$*"
	fi
}

err() {
	echo "ERROR: $*" >&2
}

if [ -z "${SSH_HOST}" ]; then
	err "Missing required environment variable: SSH_HOST"
	exit 1
fi

# Set up the host key trust
if [ ! -f "${ssh_known_hosts}" ]; then
	echo "Using ssh-keyscan to get the host key of the DDNS server"
	mkdir -p "$(dirname "${ssh_known_hosts}")"
	ssh-keyscan -t ed25519 -p "${ssh_port}" "${SSH_HOST}" | tee "${ssh_known_hosts}"
fi

cp "${ssh_known_hosts}" /etc/ssh/ssh_known_hosts

# Get old (cached) and current public IP
if [ -f "${cache_file}" ]; then
	old_ip="$(cat "${cache_file}")"
else
	old_ip=""
fi

if [ -z "$new_ip" ]; then
	# Get new IP
	new_ip="$(curl -s icanhazip.com)"
fi

ip_parser='^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$'

# Check response
if ! echo "$new_ip" | grep -q -E "${ip_parser}"; then
	err "icanhazip.com did not return an IP address"
	exit 2
fi

if [ "$old_ip" = "$new_ip" ]; then
	debug "IP has not changed."
	if [ "$force" = "yes" ]; then
		echo "Forcing update"
	else
		debug "Exiting..."
		exit
	fi
fi

echo "Updating IP from ${old_ip} to ${new_ip}"

conn="${remote_user}@${SSH_HOST}"

# The $0 of this command actually doesn't matter
# It will be overwritten by the server anyway
# Using a "command" entry in the authorized_keys file
cmd="/ddns-update.sh ${new_ip}"

ssh -i "${ssh_keyfile}" -p "${ssh_port}" "${conn}" "${cmd}"

echo "${new_ip}" >"${cache_file}"
