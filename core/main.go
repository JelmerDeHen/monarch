package main

import (
	"os"

	"github.com/JelmerDeHen/monarch/core/cmd"
)

func main() {
	run(newClient(), os.Args...)
}

func run(client *cmd.Client, args ...string) {
	app := cmd.NewApp(client)
	if err := app.Run(os.Args); err != nil {
		app.Run(args)
	}
}

func newClient() *cmd.Client {
	return &cmd.Client{}
}
