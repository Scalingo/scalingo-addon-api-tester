package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var (
	manifestPathFlag = cli.StringFlag{Name: "manifest,m", Usage: "Path to the addon manifest", Value: "manifest.json", EnvVar: "MANIFEST_PATH"}
	manifest         *Manifest
)

func main() {
	app := cli.NewApp()
	app.Name = "scalingo-addon-api-tester"
	app.Usage = "Test your addon HTTP service"
	app.Flags = []cli.Flag{manifestPathFlag}
	app.Version = "0.0.1"
	app.Author = "Leo Unbekandt"
	app.Email = "leo@scalingo.com"
	app.Before = func(c *cli.Context) error {
		var err error
		manifest, err = readManifest(c.String("manifest"))
		if err != nil {
			log.Fatalln("Failed to read addon manifest:", err)
		}
		err = manifest.Check()
		if err != nil {
			fmt.Printf("Errors in the manifest:\n%v\n", err)
			os.Exit(1)
		}
		return nil
	}
	app.Commands = []cli.Command{
		provisionCommand, deprovisionCommand, updateCommand, listCommand, purgeCommand,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("An error occured:", err)
	}
}
