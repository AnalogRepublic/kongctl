package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "kongctl"
	app.Usage = "Kong management tool"
	app.Action = func(c *cli.Context) error {
		fmt.Println("Let's manage our Kong")
		return nil
	}

	app.Run(os.Args)
}
