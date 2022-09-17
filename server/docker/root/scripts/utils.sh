#!/bin/sh

debug() {
	if [ "${DEBUG}" = 'true' ]; then
		echo "$*"
	fi
}

err() {
	>&2 echo "ERROR: $*"
}

