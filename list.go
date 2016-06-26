package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
)

var (
	listCommand = cli.Command{
		Name:  "list",
		Usage: "List saved addons",
		Action: func(c *cli.Context) {
			db, err := getDB()
			if err != nil {
				log.Fatalln("Fail to open DB:", err)
			}
			for _, entry := range *db {
				fmt.Printf("- %v: %v\n", entry.ID, entry.Plan)
			}
		},
	}
)
