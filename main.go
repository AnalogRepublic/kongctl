package main

import (
	"fmt"
	"os"

	"github.com/analogrepublic/kongctl/kong"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "kongctl"
	app.Usage = "Kong management tool"
	app.Action = func(c *cli.Context) error {
		fmt.Println("Let's manage our Kong")

		k, err := kong.NewKong("http://1270.0.0.1:8001")

		if err != nil {
			return errors.Wrap(err, "Unable to connect to Kong")
		}

		plugins, err := k.Plugins()

		if err != nil {
			return errors.Wrap(err, "Cannot fetch plugins")
		}

		return err
	}

	app.Run(os.Args)
}
