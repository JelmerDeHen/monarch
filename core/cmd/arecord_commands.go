package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/xidle"
)

func (cli *Client) Arecord(cCtx *cli.Context) error {
	args := []string{
		"-D", "sysdefault:CARD=NTUSB",
		"-t", "wav",
		"-f", "S24_3LE",
		"-r", "192000",
		"-d", "3600",
		"${OUTFILE}",
	}
	job := xidle.NewCmdJob("arecord", args...)

	job.OutfileGenerator = func() string {
		return getOutfilename("/data/mon/arecord_new/", "wav")
	}

	idlemon := xidle.NewIdlemon(job)
	idlemon.Run()

	return nil
}
