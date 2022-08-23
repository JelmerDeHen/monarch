package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/xidle"
)

func ffmpegV4l2Argv() []string {
	// Create out file
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	now := time.Now()
	outfile := fmt.Sprintf(
		"/data/mon/v4l2_new/%s.%02d%02d%02d.%02d%02d%02d.mkv",
		hostname,
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
	)

	video := "/dev/video0"
	/*
		  var video string
			name, err := findCameraByName("Logitech BRIO")
			if err != nil {
		    video = "video0"
			} else {

		  }
			video = fmt.Sprintf("/dev/%s", name)
		  fmt.Println(video)
	*/
	arguments := []string{
		"-nostdin", "-hide_banner",
		"-loglevel", "warning",
		"-f", "v4l2",
		"-an",
		//        "-input_format", "yuyv422",
		"-input_format", "mjpeg",
		"-video_size", "1920x1080",
		"-framerate", "1",
		"-i", video,
		"-vcodec", "libx265",
		outfile,
	}
	return arguments
}

// notinuse
// Finds first camera named name
// Look for name in /sys/class/video4linux/*/name
func findCameraByName(name string) (string, error) {
	///sys/class/video4linux/video0/name
	entries, err := ioutil.ReadDir("/sys/class/video4linux/")
	if err != nil {
		return "", err
	}

	if len(entries) == 0 {
		return "", fmt.Errorf("/sys/class/video4linux/ is empty")
	}

	for _, entry := range entries {
		namefn := fmt.Sprintf("/sys/class/video4linux/%s/name", entry.Name())
		_, err := os.Stat(namefn)
		if err != nil {
			continue
		}
		b, err := os.ReadFile(namefn)
		devname := strings.TrimSpace(string(b))
		if devname == name {
			return entry.Name(), nil
		}
	}
	return "", fmt.Errorf("Could not find camera named %q", name)
}

func (cli *Client) V4l2(cCtx *cli.Context) error {
	fmt.Println("Start v4l2 ", cCtx.Args().First())

	args := ffmpegV4l2Argv()
	runner := xidle.NewCmdJob("ffmpeg", args...)
	idlemon := xidle.NewIdlemon(runner)
	idlemon.Run()

	return nil
}
