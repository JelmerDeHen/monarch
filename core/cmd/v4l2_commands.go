package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/xidle"
)

func (cli *Client) V4l2(cCtx *cli.Context) error {
	args := []string{
		"-nostdin", "-hide_banner",
		"-loglevel", "warning",
		// Overwrites existing outfile
		"-y",
		// Rotate after 1 hour
		"-t", "3600",
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
		name := getOutfilename("/data/mon/v4l2", "mkv")
		return name
	}
	job.KillSignal = os.Interrupt

	idlemon := xidle.NewIdlemon(job)
	idlemon.Run()

	return nil
}
