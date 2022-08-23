package main

import (
  "os"

  "github.com/JelmerDeHen/monarch/core/cmd"
)

func main() {
	run(newProductionClient(), os.Args...)
}

func run(client *cmd.Client, args ...string) {
	app := cmd.NewApp(client)
	if err := app.Run(os.Args); err != nil {
		app.Run(args)
	}
}

func newProductionClient() *cmd.Client {
//	lggr, closeLggr := logger.NewLogger()
//	prompter := cmd.NewTerminalPrompter()
	return &cmd.Client{
	}
}
