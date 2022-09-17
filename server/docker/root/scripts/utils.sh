#!/bin/sh

debug() {
	if [ "${DEBUG}" = 'true' ]; then
		echo "$*"
	fi
}

err() {
	echo "ERROR: $*" >&2
}
