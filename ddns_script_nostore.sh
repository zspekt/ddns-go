#!/usr/bin/env sh

if [[ $# -lt 3 ]]; then
	logger DDNS missing argument
	logger DDNS proper usage: ddns_script INTERFACE ADDR PORT
	exit 1
fi

ipv4_regex="[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}"

if="$1"      # interface (wan in router)
address="$2" # addr to send the IP to
port="$3"    # port idem

get_ipv4() {
	ip addr show "$if" | grep -o "inet $ipv4_regex" | awk '{ print $2 }'
}

ipv4=$(get_ipv4)

if [[ -z "$ipv4" ]]; then
	logger DDNS ipv4 var is empty
	exit
fi

update() {
	echo "$ipv4" | ncat --send-only "$address" "$port"
}

update
