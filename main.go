package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/analogrepublic/kongctl/commands"
	"github.com/analogrepublic/kongctl/config"
	"github.com/analogrepublic/kongctl/kong"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var kongApi *kong.Kong

const (
	version     = "0.1.0"
	name        = "kongctl"
	description = "Kong management tool"
)

func main() {
	var err error

	err = config.Init()

	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to read config file"))
		os.Exit(1)
	}

	c := config.GetConfig()

	app := cli.NewApp()

	app.Name = name
	app.Usage = description
	app.Version = version
	app.EnableBashCompletion = true

	context, err := c.GetCurrentContext()

	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to determine context to use, please check your config file"))
		os.Exit(1)
	}

	kongApi, err = kong.NewKong(context.Host, c)

	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to communicate with the Kong service"))
		os.Exit(1)
	}

	commands.SetKongApi(kongApi)

	app.Commands = []cli.Command{
		commands.Context,
		commands.Apply,
		commands.Export,
		commands.Describe,
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}
