#!/usr/bin/env bash
# Looks for unused files in directory specified by indir and compresses them with ffmpeg to compress them to outdir
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
. "${DIR}/SYSTEMD_LSOF.sh"

read -r -d '' CONFIG_SREC <<'EOF'
{
	"service": "mon_srec",
	"indir": "/data/mon/srec",
	"outdir": "/data/mon/srec_mp4",
	"trim": ".mkv",
	"ext": ".mp4",
	"argv": [ "-an" ],
	"dryrun": false
}
EOF
read -r -d '' CONFIG_ARECORD <<'EOF'
{
	"service": "mon_arecord",
	"indir": "/data/mon/arecord",
	"outdir": "/data/mon/arecord_mp3",
	"trim": ".wav",
	"ext": ".mp3",
	"argv": [ ],
	"dryrun": false
}
EOF
read -r -d '' CONFIG_VIDEO <<'EOF'
{
	"service": "mon_video",
	"indir": "/data/mon/video",
	"outdir": "/data/mon/video_mp4",
	"trim": ".mkv",
	"ext": ".mp4",
	"argv": [ "-an" ],
	"dryrun": false
}
EOF
read -r -d '' CONFIG_V4L2 <<'EOF'
{
	"service": "mon_v4l2",
	"indir": "/data/mon/v4l2",
	"outdir": "/data/mon/v4l2_mp4",
	"trim": ".mkv",
	"ext": ".mp4",
	"argv": [ "-an" ],
	"dryrun": false
}
EOF
# A recording can be deleted after the file has been compressed
function rmProcessed () {
	local _IN="${1}" _OUT="${2}"
	#if [ ! -s "${_OUT}" ]; then
	#	printf '%s: Need nonempty file at %s\n' "${FUNCNAME[0]}" "${_OUT}"
	#	return 1
	#fi
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
function compress_ffmpeg () {
	local _IN="${1}" _OUT="${2}" _ARGV="${3}"
	# Collect all args after arg 2 and put them in array for exec
	shift 2 ; printf -v _ARGV '%s ' "${@}"
	# Sanity checks
	if [ ! -f "${_IN}" ]; then
		printf '%s: Need input file %s\n' "${FUNCNAME[0]}" "${_IN}"
		return 1
	fi
	if [ ! -s "${_IN}" ]; then
		printf '%s: Need nonempty input file %s\n' "${FUNCNAME[0]}" "${_IN}"
		return 1
	fi
	if [ -f "${_OUT}" ]; then
		printf '%s: %s already exists\n' "${FUNCNAME[0]}" "${_IN}"
		return 1
	fi
	local -a REPLIES
	while read -r; do
		REPLIES+=("${REPLY}")
		printf '%s\n' "${REPLY}"
	done < <( /usr/bin/env -i /usr/bin/ffmpeg </dev/null -loglevel warning -i "${_IN}" ${_ARGV[@]} "${_OUT}" 2>&1 )
	if [ "${?}" == "0" ]; then
		printf '%s\n' "${REPLIES[@]}" > "${_OUT}.status"
	fi
}
# getOutputFileWithInputFile <in> <trimFromSuffix> <extension> <outdir>
# getOutputFileWithInputFile /tmp/in/test.wav .wav .mp3 /tmp/out
function getOutputFileWithInputFile  () {
	local _IN="${1}" _TRIMSUFFIX="${2}" _EXTENSION="${3}" _OUTDIR="${4}" _OUT=
	_OUT="${_IN%${_TRIMSUFFIX}}${_EXTENSION}"
	printf '%s\n' "${_OUTDIR}/${_OUT##*/}"
}
function compress() {
	if [ -z "${1}" ]; then
		printf '%s: Need json CONFIG as first argument' "${FUNCNAME[0]}"
		return
	fi
	local _CONFIG="${1}"
	# Fetch conf vars
	local _SERVICE=$(jq -r .service <<<"${_CONFIG}")
	local _INDIR=$(jq -r .indir <<<"${_CONFIG}")
	local _OUTDIR=$(jq -r .outdir <<<"${_CONFIG}")
	local _TRIM=$(jq -r .trim <<<"${_CONFIG}")
	local _EXT=$(jq -r .ext <<<"${_CONFIG}")
	local _DRYRUN=$(jq -r '.dryrun | if . != "true" then . |=  "false" else . end ' <<<"${_CONFIG}")
	local _ARGV=$(jq -r '.argv | .[]' <<<"${_CONFIG}")
	echo "${_ARGV}"

	# Sanity checks
	local _ERR=0
#	if ! systemctl is-active --quiet "${_SERVICE}" &>/dev/null; then
#		printf '%s: Service %s is not active\n' "${FUNCNAME[0]}" "${_SERVICE}"
#		let _ERR+=1
#	fi
	if [ ! -d "${_INDIR}" ]; then
		printf '%s: input directory %s does not exist\n' "${FUNCNAME[0]}" "${_INDIR}"
		let _ERR+=1
	fi
	if [ ! -d "${_OUTDIR}" ]; then
		printf '%s: outdir %s does not exist\n' "${FUNCNAME[0]}" "${_OUTDIR}"
		if ! mkdir -pv "${_OUTDIR}"; then
			let _ERR+=1
		fi
	fi
	local _TMPFILE="${_OUTDIR}/$RANDOM.tmp"
	touch "${_TMPFILE}"
	if [ $? -ne 0 ]; then
		printf '%s: outdir %s is not writable\n' "${FUNCNAME[0]}" "${_OUTDIR}"
		let _ERR+=1
	fi
	rm -f "${_TMPFILE}"
	if [ "${_ERR}" -ne 0 ]; then
		printf '%s: %d errors to fix\n' "${FUNCNAME[0]}" "${_ERR}"
		return 1
	fi
	# Crawl indir
	local _LOSSYFILE
	while read -r _LOSSYFILE; do
		# Check if file is in-use
		if SYSTEMD_SERVICE_OPENED_FILE "${_SERVICE}" "${_LOSSYFILE}"; then
			printf '%s: Service %s is using %s\n' "${FUNCNAME[0]}" "${_SERVICE}" "${_LOSSYFILE}"
			continue
		fi
		printf '%s: Compressing %s\n' "${FUNCNAME[0]}" "${_LOSSYFILE}"
		# Build output file based on input file and config
		local _COMPRESSEDOUT=$( getOutputFileWithInputFile "${_LOSSYFILE}" "${_TRIM}" "${_EXT}" "${_OUTDIR}" )
		if [ -f "${_COMPRESSEDOUT}" ]; then
			# Does the .status file exist for this output
			# This file is only generated when the compression was completed and had exit status 0
			if [ -f "${_COMPRESSEDOUT}.status" ]; then
				if [ "${_DRYRUN}" != "true" ]; then
					rmProcessed "${_LOSSYFILE}" "${_COMPRESSEDOUT}"
				fi
			else
				printf '%s: Found compressed file %s but not %s.status signaling the compression process failed previously\n' "${FUNCNAME[0]}" "${_COMPRESSEDOUT}" "${_COMPRESSEDOUT}"
				rm -v "${_COMPRESSEDOUT}"
			fi
		fi
		# Remove empty input files
		[ -s "${_LOSSYFILE}" ] || rmProcessed "${_LOSSYFILE}" "${_COMPRESSEDOUT}"
		compress_ffmpeg "${_LOSSYFILE}" "${_COMPRESSEDOUT}" "${_ARGV}"
		#if [ "$?" != "0" ]; then
		#	printf '%s: Error processing %s\n' ""
		#fi
	done < <( find "${_INDIR}" -type f )
}
case "${1^^}" in
	ARECORD)
		compress "${CONFIG_ARECORD}"
		;;
	SREC)
		compress "${CONFIG_SREC}"
		;;
	VIDEO)
		compress "${CONFIG_VIDEO}"
		;;
	V4L2)
		compress "${CONFIG_V4L2}"
		;;
	*)
		printf 'No available for CONFIG profile %s\n' "${1}"
		;;
esac
