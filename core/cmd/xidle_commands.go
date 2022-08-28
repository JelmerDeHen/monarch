package cmd

import (
	"time"

	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/xidle"
)

func (cli *Client) Xidle(cCtx *cli.Context) error {
	job := xidle.NewCmdJob("sleep", "1337")

	// Configure durations to something we can wait for
	// When user present last second spawn
	// When user idle over 5 secs kill
	idlemon := xidle.NewIdlemon(job)
	idlemon.IdleLessTimeout = time.Second
	idlemon.IdleOverTimeout = time.Second * 5

	idlemon.Run()

	return nil
}
