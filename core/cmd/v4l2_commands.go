package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/xidle"
)

func (cli *Client) V4l2(cCtx *cli.Context) error {
	args := []string{
		"-nostdin", "-hide_banner",
		"-loglevel", "warning",
    "-y",
		"-f", "v4l2",
		"-an",
		//        "-input_format", "yuyv422",
		"-input_format", "mjpeg",
		"-video_size", "1920x1080",
		"-framerate", "1",
		"-i", "/dev/video0",
		"-vcodec", "libx265",
		"${OUTFILE}",
	}
	job := xidle.NewCmdJob("ffmpeg", args...)

	job.OutfileGenerator = func() string {
		return getOutfilename("/data/mon/v4l2_new/", "mkv")
	}

	idlemon := xidle.NewIdlemon(job)
	idlemon.Run()

	return nil
}
