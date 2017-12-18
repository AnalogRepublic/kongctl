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

		_, err := kong.NewKong("http://testing.kongctl.io:8001")

		if err != nil {
			return errors.Wrap(err, "Unable to connect to Kong")
		}

		return nil
	}

	app.Run(os.Args)
}
