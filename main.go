package main

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "vault-cli"
	app.Usage = ""

	app.Action = func(c *cli.Context) error {
		conf := api.DefaultConfig()
		client, _ := api.NewClient(conf)

		client.SetToken("64f37837-261b-7333-9b0b-35c2cc3e7e17")
		paths, err := unroll("secret/", client)
		if err != nil {
			return err
		}

		search := c.Args().Get(0)
		if strings.TrimSpace(search) == "" {
			return errors.New("Please provide a search term")
		}

		rex := regexp.MustCompile(`.*` + search + `.*`)
		for _, path := range paths {
			if rex.MatchString(path) {
				res, err := client.Logical().Read(path)
				if err != nil || res == nil {
					continue
				}

				data := res.Data
				data["path"] = path

				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				enc.Encode(data)
			}
		}
		return nil
	}

	app.Run(os.Args)
}
