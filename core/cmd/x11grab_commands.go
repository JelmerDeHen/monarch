package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/scrnsaver"
	"github.com/JelmerDeHen/xidle"
)

func (cli *Client) X11grab(cCtx *cli.Context) error {
	arguments := []string{
		"-nostdin", "-hide_banner",
		"-loglevel", "warning",
		"-y",
		"-f", "x11grab",
		"-an",
		"-framerate", "25",
		"-video_size", scrnsaver.GetResolution(),
		"-i", os.Getenv("DISPLAY"),
		"-vcodec", "libx265",
		"-preset", "ultrafast",
		"${OUTFILE}",
	}

	job := xidle.NewCmdJob("ffmpeg", arguments...)
	job.OutfileGenerator = func() string {
		return getOutfilename("/data/mon/srec_new/", "mkv")
	}

	idlemon := xidle.NewIdlemon(job)
	idlemon.Run()

	return nil
}
