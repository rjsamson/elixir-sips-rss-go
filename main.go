package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/rjsamson/elixir-sips-rss-go/downloader"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Elixir Sips RSS Downloader"
	app.Usage = "Usage information"
	app.Version = "0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "username, u",
			Usage: "elixir sips username",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "elixir sips password",
		},
		cli.IntFlag{
			Name:  "episodes, e",
			Usage: "number of episodes to download (default 5)",
			Value: 5,
		},
	}

	app.Action = func(c *cli.Context) {
		episodes := c.Int("episodes")
		username := c.String("username")
		password := c.String("password")

		if username == "" || password == "" {
			fmt.Println("Please enter a username and password")
			return
		}

		downloader.Download(username, password, episodes)
	}

	app.Run(os.Args)
}
