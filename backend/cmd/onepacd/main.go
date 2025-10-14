package main

import (
	"log"
	"os"

	"github.com/frimin/1pactus-react/backend/app/onepacd"
	"github.com/urfave/cli/v2"
)

func main() {
	cmd := cli.NewApp()
	cmd.Name = onepacd.App
	cmd.Version = onepacd.Version
	cmd.Commands = []*cli.Command{
		{
			Name: "run", Usage: "App Run",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:     "config",
					Aliases:  []string{"c"},
					Usage:    "Load configuration from `FILE`",
					Required: true,
				},
				&cli.StringSliceFlag{
					Name:     "param",
					Aliases:  []string{"p"},
					Usage:    "configuration cli overrides",
					Required: false,
				},
			},
			Action: func(c *cli.Context) error {
				if err := onepacd.LoadConfig(onepacd.App, c.StringSlice("config"), c.StringSlice("param")); err != nil {
					return err
				}
				onepacd.Run()
				return nil
			},
		},
	}
	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
