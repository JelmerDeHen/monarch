#!/bin/sh
# Capture from virtual x11grab virtual device
# https://ffmpeg.org/ffmpeg.html#Generic-options
# http://www.ffmpeg.org/ffmpeg-devices.html#x11grab
function ffmpeg_x11grab () {
	local _DIMENSIONS _OUTFILE
	_OUTFILE="${1}"
	DISPLAY=":0"
	# Get current dimensions
	_DIMENSIONS="$( DISPLAY=${DISPLAY} xdpyinfo | grep -om1 '[0-9]\{4\}x[0-9]\{4\}' ; )"
	if [ -z "${_OUTFILE}" ]; then
		printf >&2 'Need outfile!\n'
		return 1
	elif [ -z "${_DIMENSIONS}" ]; then
		printf >&2 'Need dimensions!\n'
		return 1
	fi
	if ! touch "${_OUTFILE}"; then
		printf >&2 '%s: Could not write to output file "%s"\n' "${FUNCNAME[0]}" "${_OUTFILE}"
		return 1
	fi
	rm "${_OUTFILE}"
	# -an blocks audio streams
	# -framerate
	printf 'DIMENSIONS=%s; OUTFILE=%s\n' "${_DIMENSIONS}" "${_OUTFILE}"
	ffmpeg \
		-timelimit 3600\
		-loglevel warning\
		-nostdin\
		-f x11grab\
		-an\
		-r 25\
		-video_size "${_DIMENSIONS}"\
		-i "${DISPLAY}" \
		-c:v libx264\
		-preset ultrafast\
		"${_OUTFILE}"
	return $?
}
ffmpeg_x11grab "${@}"
exit $?

