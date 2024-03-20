package main

import (
	"github.com/apex/log"
	"os"

	"github.com/urfave/cli/v2"
	bootstrap "github.com/victorsantoso/endeus/cmd/cli"
)

var commands = []*cli.Command{
	{
		Name: "start", 
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "config",
				Aliases: []string{"c"},
				Usage: "-c path will be used for config eg: -c ./config.json",
			},
		},
		Action: func(ctx *cli.Context) error {
			config := ctx.String("config")
			return bootstrap.Bootstrap(config)
		},
	},
}

func main() {
	app := &cli.App{
		Name: "endeus",
		Version: "1.0.1",
		Commands: commands,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("[main] error running application, err: %v", err)
	}
}
