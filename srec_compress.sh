#!/usr/bin/env bash
# Look for unused files in _INDIR and format them to 
_INDIR=/data3/mon/srec
_OUTDIR=/data3/mon/srec_mp4
function compressVideo () {
	local _IN="${1}" _OUT="${2}"
	if [ ! -s "${_IN}" ]; then
		printf '%s: Need nonempty input file %s\n' "${FUNCNAME[0]}" "${_IN}"
		return 1
	fi
	if [ -f "${_OUT}" ]; then
		printf '%s: %s already exists\n' "${FUNCNAME[0]}" "${_IN}"
		return 1
	fi
	ffmpeg -i "${_IN}" -an "${_OUT}"
}
# A recording can be deleted after the file has been compressed
function rmProcessed () {
	local _IN="${1}" _OUT="${2}"
	if [ ! -s "${_IN}" ]; then
		printf '%s: Need nonempty file at %s\n' "${FUNCNAME[0]}" "${_IN}"
		return 1
	fi
	if [ ! -s "${_OUT}" ]; then
		printf '%s: Need nonempty file at %s\n' "${FUNCNAME[0]}" "${_OUT}"
		return 1
	fi
	if lsof "${_IN}" &>/dev/null; then
		printf '%s: %s is busy\n' "${FUNCNAME[0]}" "${_IN}"
		return 1
	fi
	if lsof "${_OUT}" &>/dev/null; then
		printf '%s: %s is busy\n' "${FUNCNAME[0]}" "${_OUT}"
		return 1
	fi
	rm -v "${_IN}"
}
function main() {
	local _MP4OUT
	while read -r; do
		# Test if not in-use
		if lsof "${REPLY}" &>/dev/null; then
			printf '%s: %s is busy!\n' "${FUNCNAME[0]}" "${REPLY}"
			continue
		fi
		_MP4OUT="${REPLY%.*}.mp4" _MP4OUT="${_OUTDIR}/${_MP4OUT##*/}"
		if [ -f "${_MP4OUT}" ]; then
			rmProcessed "${REPLY}" "${_MP4OUT}"
			continue
		fi
		compressVideo "${REPLY}" "${_MP4OUT}"
	done < <( find "${_INDIR}" -type f )
}
while :; do
	main
	printf 'Sleeping 2 mins\n'
	sleep 120
done
