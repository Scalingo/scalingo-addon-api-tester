package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/Scalingo/scalingo-addon-api-tester/Godeps/_workspace/src/github.com/codegangsta/cli"
)

var (
	provisionCommand = cli.Command{
		Name:  "provision",
		Usage: "Request addon provisioning",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "plan", Usage: "Optional - Choose the plan to provision, choose the first of the manifest if empty"},
			cli.StringFlag{Name: "app", Usage: "Optional - Choose a name, will generate one if not defined"},
		},
		Action: func(c *cli.Context) {
			plan := c.String("plan")
			if plan == "" {
				plan = manifest.Plans[0].Name
			} else if !manifest.PlanExist(plan) {
				log.Fatalln("Plan", plan, "is not defined in manifest")
			}
			options := manifest.PlanOptions(plan)
			app := c.String("app")
			if app == "" {
				app = generateName()
			}
			res, err := doRequest("POST", manifest.Test.BaseURL, map[string]interface{}{
				"plan":    plan,
				"app_id":  app,
				"options": options,
			})
			if err != nil {
				log.Fatalln("Fail to provision resource:", err)
			}
			defer res.Body.Close()
			if res.StatusCode != 200 && res.StatusCode != 201 && res.StatusCode != 202 {
				log.Fatalln("Addon returned bad status:", res.Status, "expected 201, 202 (or exceptionnaly 200 for heroku compatibility)")
			}
			var proRes ProvisioningResponse
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalln("Failed to read body from addon:", err)
			}
			err = json.Unmarshal(body, &proRes)
			if err != nil {
				log.Fatalln("Addon returned bad JSON:", err, "-", string(body))
			}
			if proRes.ID == "" {
				log.Fatalln("Addon returned an empty ID")
			}
			if proRes.Message == "" {
				log.Println("Notice: no message received")
			}
			if proRes.Config == nil || len(proRes.Config) == 0 {
				log.Println("Notice: no configuration received")
			}
			if err = manifest.CheckAddonConfig(proRes.Config); err != nil {
				log.Fatalln(err)
			}
			err = saveAddonRef(proRes.ID, plan)
			if err != nil {
				log.Println("fail to save addon reference:", err)
			}
			fmt.Println("→ OK")
		},
	}
)

var (
	names = []string{"african-wild-dog", "amur-leopard", "amur-tiger", "asian-elephant", "bengal-tiger", "black-rhino", "black-spider-monkey", "black-footed-ferret", "blue-whale", "bluefin-tuna", "bonobo", "bornean-orangutan", "borneo-pygmy-elephant", "chimpanzee", "cross-river-gorilla", "eastern-lowland-gorilla", "endangered-species", "fin-whale", "galápagos-penguin", "ganges-river-dolphin", "giant-panda", "green-turtle", "hawksbill-turtle", "humphead-wrasse", "indian-elephant", "indochinese-tiger", "indus-river-dolphin", "javan-rhino", "leatherback-turtle", "loggerhead-turtle", "malayan-tiger", "mountain-gorilla", "north-atlantic-right-whale", "orangutan", "our-work-protecting-species", "pangolin", "saola", "sea-lions", "sei-whale", "snow-leopard", "south-china-tiger", "species", "sri-lankan-elephant", "sumatran-elephant-", "sumatran-orangutan", "sumatran-rhino", "sumatran-tiger", "tiger", "vaquita", "western-lowland-gorilla", "whale", "yangtze-finless-porpoise"}
)

func generateName() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(names))
	n := rand.Intn(10000)
	return fmt.Sprintf("%v-%v", names[i], n)
}
