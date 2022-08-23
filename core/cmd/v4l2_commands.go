package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/monarch/core/v4l2"
)

func (cli *Client) V4l2 (cCtx *cli.Context) error {
	fmt.Println("Start x11grab ", cCtx.Args().First())
	v4l2.Run()
	return nil
}
