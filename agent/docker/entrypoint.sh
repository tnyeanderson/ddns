#!/bin/sh

### CONFIGURATION: Set these as environment variables
#SSH_HOST="bar.com"
#SSH_PORT="2222"

# Set value to "force" to make update the DDNS entry
# even if it hasn't changed
force=$1

remote_user=ddns
cache_file=/ddns/current-ip
ssh_keyfile=/ddns/ssh.key

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

# Get old (cached) and current public IP
if [ -f "${cache_file}" ]; then
	old_ip="$(cat "${cache_file}")"
else
	old_ip=""
fi

# Get new IP
new_ip="$(curl -s icanhazip.com)"

ip_parser='^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$'

# Check response
if ! echo "$new_ip" | grep -q -E "${ip_parser}"; then
	err "icanhazip.com did not return an IP address"
	exit 2
fi

if [ "$old_ip" = "$new_ip" ]; then
	debug "IP has not changed."
	if [ "$force" = "force" ]; then
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

ssh -i "${ssh_keyfile}" -o "StrictHostKeyChecking no" -p "${SSH_PORT}" "${conn}" "${cmd}"

echo "${new_ip}" >"${cache_file}"
