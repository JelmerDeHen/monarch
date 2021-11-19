#!/bin/bash

function SYSTEMD_LSOF () {
	local _SERVICE="${1}"
	if ! systemctl is-active --quiet "${_SERVICE}" &>/dev/null; then
		printf >&2 '%s: %s is not active!\n' "${FUNCNAME[0]}" "${_SERVICE}"
		return 1
	fi
	local _MAINPID="$(systemctl show --property MainPID --value -- ${_SERVICE})"
	local _PID _OPENFILE _STATUSFILE
	while read -r _STATUSFILE; do
		_PID="${_STATUSFILE//[^0-9]/}"
		#printf '# PID=%d (%s)\n' "${_PID}" "$(</proc/${_PID}/comm)"
		while read -r _OPENFILE; do
			if [ "${_OPENFILE#/}" = "${_OPENFILE}" ] || \
				[ "${_OPENFILE#/proc}" != "${_OPENFILE}" ] || \
				[ "${_OPENFILE#/dev}" != "${_OPENFILE}" ]; then
				continue
			fi
			_OPENFILE="${_OPENFILE% (deleted)}"
			printf '%s\n' "${_OPENFILE}"
		done < <(find /proc/${_PID}/fd -type l -exec readlink -f "{}" \; )
	done < <( printf '/proc/%d/status\n' "${_MAINPID}"; grep -sl $'^PPid:\t'"${_MAINPID}"'$' /proc/[0-9]*/status )
	#if [ "${_OK}" = "0" ]; then
	#	printf >&2 '%s: MainPID of service %s (%d) had no child procs to enumerate\n' "${FUNCNAME[0]}" "${_SERVICE}" "${_MAINPID}"
	#fi
	return 0
}
function SYSTEMD_OPENED () {
	local _SERVICE
	while read -r _SERVICE REPLY; do
		[ "${_SERVICE%.service}" = "${_SERVICE}" ] && continue
		echo "# SERVICE=${_SERVICE}"
		SYSTEMD_LSOF "${_SERVICE}"
	done < <( systemctl list-units  --type service --state running )
}

function SYSTEMD_SERVICE_OPENED_FILE () {
	local _SERVICE="${1}" _FILE="${2}" _RET
	while read -r; do
		if [ "${REPLY}" = "${_FILE}" ]; then
			return 0
		fi
	done < <( SYSTEMD_LSOF "${_SERVICE}"; )
	return 1
}
#SYSTEMD_SERVICE_OPENED_FILE mon_arecord /data3/mon/arecord/desktop.20200706.1639.wav
#SYSTEMD_OPENED
#SYSTEMD_LSOF "${1}"
