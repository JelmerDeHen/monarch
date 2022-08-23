package cmd

import (
	"fmt"
	//  "time"

	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/xidle"
)

func (cli *Client) Xidle(cCtx *cli.Context) error {
	fmt.Println("Start idlecmd ", cCtx.Args().First())

	runner := xidle.NewCmdJob("sleep", "1337")
	idlemon := xidle.NewIdlemon(runner)
	idlemon.Run()

	return nil
}
