package cmd

import (
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/xidle"
)

func (cli *Client) X11grab(cCtx *cli.Context) error {
	arguments := []string{
		"-nostdin", "-hide_banner",
		"-loglevel", "warning",
		"-y",
		"-t", "3600",
		"-f", "x11grab",
		"-an",
		"-r", "1",
		"-video_size", "${RESOLUTION}",
		"-i", os.Getenv("DISPLAY"),
		"-vcodec", "libx265",
		"-preset", "ultrafast",
		"${OUTFILE}",
	}

	job := xidle.NewCmdJob("ffmpeg", arguments...)
	job.OutfileGenerator = func() string {
		return getOutfilename("/data/mon/x11grab", "mkv")
	}
	job.KillSignal = os.Interrupt

	idlemon := xidle.NewIdlemon(job)
	//idlemon.IdleOverT = time.Second * 5
	//idlemon.IdleLessT = time.Second

	idlemon.Run()

	return nil
}
