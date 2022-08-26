package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/JelmerDeHen/monarch/core/compress"
)

func lsPipelines() []string {
	var pipelines []string
	for k, _ := range compress.Pipelines {
		pipelines = append(pipelines, k)
	}
	return pipelines

}

func (cli *Client) Compress(cCtx *cli.Context) error {
	// When a pipeline is provided scan that pipeline
	pipeline := cCtx.Args().First()
	if pipeline != "" {
		if _, ok := compress.Pipelines[pipeline]; !ok {
			pipelines := lsPipelines()
			fmt.Printf("Pipeline %s not found, the following pipelines exist: %s\n", pipeline, strings.Join(pipelines[:], ", "))
			return nil
		}

		err := compress.Pipelines[pipeline].Scan()
		if err != nil {
			fmt.Println(err)
		}
		return err
	}

	// When no pipeline option is provided scan all the pipelines
	for _, p := range compress.Pipelines {
		err := p.Scan()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	return nil
}
