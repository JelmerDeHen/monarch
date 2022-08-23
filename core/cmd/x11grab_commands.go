package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/scrnsaver"
	"github.com/JelmerDeHen/xidle"
)

func ffmpegX11grabArgv() []string {
	// Create out file
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	now := time.Now()
	outfile := fmt.Sprintf(
		"/data/mon/srec_new/%s.%02d%02d%02d.%02d%02d%02d.mkv",
		hostname,
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
	)

	display := os.Getenv("DISPLAY")
	if display == "" {
		panic("$DISPLAY was empty")
	}

	arguments := []string{
		"-nostdin", "-hide_banner",
		"-loglevel", "warning",
		"-f", "x11grab",
		"-an",
		"-framerate", "25",
		"-video_size", scrnsaver.GetResolution(),
		"-i", os.Getenv("DISPLAY"),
		"-vcodec", "libx265",
		"-preset", "ultrafast",
		outfile,
	}
	return arguments
}

func (cli *Client) X11grab(cCtx *cli.Context) error {
	fmt.Println("Start x11grab ", cCtx.Args().First())

	args := ffmpegX11grabArgv()
	runner := xidle.NewCmdJob("ffmpeg", args...)
	idlemon := xidle.NewIdlemon(runner)
	idlemon.Run()

	return nil
}
