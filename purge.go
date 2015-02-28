package main

import (
	"fmt"
	"log"

	"github.com/Scalingo/scalingo-addon-api-tester/Godeps/_workspace/src/github.com/codegangsta/cli"
)

var (
	purgeCommand = cli.Command{
		Name:  "purge",
		Usage: "Delete all the saved addons",
		Action: func(c *cli.Context) {
			err := rmDB()
			if err != nil {
				log.Fatalln("Fail to remove DB:", err)
			}
			fmt.Println("→ OK")
		},
	}
)
