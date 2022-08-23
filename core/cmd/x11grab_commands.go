package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/monarch/core/x11grab"
)

func (cli *Client) X11grab (cCtx *cli.Context) error {
	fmt.Println("Start x11grab ", cCtx.Args().First())

	x11grab.Run()
	return nil
}
