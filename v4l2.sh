#!/bin/sh
# Capture from webcam
function ffmpeg_v4l2 () {
	local _OUTFILE _VIDEO
	_VIDEO=/dev/video0
	_OUTFILE="${1}"
	if [ -z "${_OUTFILE}" ]; then
		printf >&2 'Need outfile!\n'
		return 1
	fi
	if ! touch "${_OUTFILE}"; then
		printf >&2 '%s: Could not write to output file "%s"\n' "${FUNCNAME[0]}" "${_OUTFILE}"
		return 1
	fi
	rm "${_OUTFILE}"
	printf 'OUTFILE=%s\n' "${_OUTFILE}"
	ffmpeg \
		-timelimit 3600\
		-loglevel warning\
		-nostdin\
		-f v4l2\
		-framerate 1\
		-video_size 4096x2160\
		-i ${_VIDEO}\
		"${_OUTFILE}"
	return $?
}
ffmpeg_v4l2 "${@}"
exit $?

