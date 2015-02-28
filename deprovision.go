package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Scalingo/scalingo-addon-api-tester/Godeps/_workspace/src/github.com/codegangsta/cli"
)

var (
	deprovisionCommand = cli.Command{
		Name:  "deprovision",
		Usage: "Request addon deprovisioning",
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				log.Fatalln(os.Args, "<addon id>")
			}
			id := c.Args()[0]
			res, err := doRequest("DELETE", manifest.Test.BaseURL+"/"+id, nil)
			if err != nil {
				log.Fatalln("Fail to delete resource:", err)
			}
			defer res.Body.Close()
			if res.StatusCode != 204 {
				log.Fatalln("Addon returned bad status:", res.Status, "expected 204")
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalln("Failed to read body from addon:", err)
			}
			if len(body) != 0 {
				log.Fatalln("Expected empty body, got", string(body))
			}
			err = deleteAddonRef(id)
			if err != nil {
				log.Println("fail to delete addon reference:", err)
			}
			fmt.Println("→ OK")
		},
	}
)
