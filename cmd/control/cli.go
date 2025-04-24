package control

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func Main() {
	//log.InitLogger()
	cli.VersionPrinter = func(c *cli.Context) {
		runningName := "ContOS 1.0"
		fmt.Fprintf(c.App.Writer, "version %s from os image %s\n", c.App.Version, runningName)
	}
	app := cli.NewApp()

	app.Name = os.Args[0]
	app.Usage = fmt.Sprintf("Control and configure ContOS\nbuilt: %s", "01.01.01")
	app.Version = "1.0"
	app.Author = "Vsevolod Chalkov"
	app.EnableBashCompletion = true
	app.Before = func(c *cli.Context) error {
		if os.Geteuid() != 0 {
			//log.Fatalf("%s: Need to be root", os.Args[0])
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:            "entrypoint",
			Hidden:          true,
			HideHelp:        true,
			SkipFlagParsing: true,
			Action:          entrypointAction,
		},
	}

	app.Run(os.Args)
}
