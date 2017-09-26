package main

import (
	"fmt"
	"os"

	"github.com/analogrepublic/kongctl/kong"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "kongctl"
	app.Usage = "Kong management tool"
	app.Action = func(c *cli.Context) error {
		fmt.Println("Let's manage our Kong")

		errors := make(chan error)

		go func() {
			k := kong.NewKong("http://1270.0.0.1:8001", errors)
			plugins, err := k.Plugins()
		}()

		select {
		case err := <-errors:
			fmt.Println(err)
		}

		return nil
	}

	app.Run(os.Args)
}
