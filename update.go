package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var (
	updateCommand = cli.Command{
		Name:  "update",
		Usage: "Request addon plan change",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "plan", Usage: "Choose the plan to change to"},
		},
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				fmt.Println("Usage:", os.Args[0], "<addon id> --plan <plan>")
				os.Exit(-127)
			}
			id := c.Args()[0]
			plan := c.String("plan")
			if plan == "" {
				log.Fatalln("no plan specified")
			} else if !manifest.PlanExist(plan) {
				fmt.Println("Plan", plan, "is not defined in manifest")
				os.Exit(1)
			}
			options := manifest.PlanOptions(plan)
			res, err := doRequest("PUT", manifest.Test.BaseURL+"/"+id, map[string]interface{}{
				"plan":    plan,
				"options": options,
			})
			if err != nil {
				log.Fatalln("Fail to change plan of addon:", err)
			}
			defer res.Body.Close()
			if res.StatusCode != 200 {
				log.Fatalln("Addon returned bad status:", res.Status, "expected 200")
			}
			var updRes UpdateResponse
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalln("Failed to read body from addon:", err)
			}
			err = json.Unmarshal(body, &updRes)
			if err != nil {
				log.Fatalln("Addon returned bad JSON:", err, "-", string(body))
			}
			if updRes.Message == "" {
				log.Println("Notice: no message received")
			}
			if updRes.Config == nil || len(updRes.Config) == 0 {
				log.Println("Notice: no configuration received")
			}
			if err = manifest.CheckAddonConfig(updRes.Config); err != nil {
				log.Fatalln(err)
			}
			err = saveAddonRef(id, plan)
			if err != nil {
				log.Println("fail to save addon reference:", err)
			}
			fmt.Println("→ OK")
		},
	}
)
