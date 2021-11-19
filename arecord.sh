#!/bin/sh
#[Service]
#ExecStartPre=-killall arecord
#ExecStart=sh -c 'arecord -D sysdefault:CARD=NTUSB -t wav -f S24_3LE -r 192000 -d 3600 /data3/mon/arecord/desktop.$$({date,+%%Y%%m%%d.%%H%%M}).wav'
#ExecStop=-killall arecord
#Restart=always
#RestartSec=1
#[Install]
#WantedBy=multi-user.target

# https://alsa.opensrc.org/PCM
# PCM=Pulse Code Modulation


function arecordwrap () {
	local _OUTFILE _PCM
	# _PID _PIDFILE
	_OUTFILE="/data3/mon/arecord/desktop.$({date,+%Y%m%d.%H%M}).wav"
	#_PIDFILE="/run/mon_arecord.pid"
	#if [ -s "${_PIDFILE}" ]; then
	#	printf -v _PID %d "$(<${_PIDFILE})" || return 1
	#	printf '%s: %s pid=%d\n' "${FUNCNAME[0]}" "${_PIDFILE}" "${_PID}"
	#	if [ -d "/proc/${_PID}" ]; then
	#		printf '%s: kill %d\n' "${FUNCNAME[0]}" "${_PID}"
	#		kill "${_PID}" || return 1
	#	fi
	#	printf '%s: Truncating %s!\n' "${_FUNCNAME[0]}" "${_PIDFILE}"
	#	:>"${_PIDFILE}"
	#fi
	if ! touch "${_OUTFILE}"; then
		printf >&2 '%s: Could not write "%s"!' "${FUNCNAME[0]}" "${_OUTFILE}"
		return 1
	fi
	_PCM="sysdefault:CARD=NTUSB"
	printf '_OUTFILE=%s; _PCM=%s\n' "${_OUTFILE}" "${_PCM}"
	arecord\
		-D sysdefault:CARD=NTUSB\
		-t wav\
		-f S24_3LE\
		-r 192000\
		-d 3600\
		"${_OUTFILE}"
	#printf >"${_PIDFILE}" '%d' "${!}"
}	
arecordwrap
exit $?
