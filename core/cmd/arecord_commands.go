package cmd

import (
  "fmt"
  "os"
  "time"

	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/xidle"
)

func arecordArgv() []string {
	// Create out file
  hostname, err := os.Hostname()
	if err != nil {
    panic(err)
	}
	now := time.Now()
	outfile := fmt.Sprintf(
		"/data/mon/arecord_new/%s.%02d%02d%02d.%02d%02d%02d.wav",
		hostname,
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
	)

	arguments := []string{
		"-D", "sysdefault:CARD=NTUSB",
    "-t", "wav",
    "-f", "S24_3LE",
    "-r", "192000",
    "-d", "3600",
    outfile,
	}
	return arguments
}

func (cli *Client) Arecord (cCtx *cli.Context) error {
	fmt.Println("Start arecord ", cCtx.Args().First())

  args := arecordArgv()
  runner := xidle.NewCmdJob("arecord", args...)
  idlemon := xidle.NewIdlemon(runner)
	idlemon.Run()

  return nil
}
