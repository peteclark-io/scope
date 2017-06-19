package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "scope"
	app.Usage = "Scopes out vault keys and reads the first match"

	app.Action = func(c *cli.Context) error {
		conf := api.DefaultConfig()
		client, _ := api.NewClient(conf)

		err := setToken(client)
		if err != nil {
			return err
		}

		paths, err := unroll("secret/", client)
		if err != nil {
			return err
		}

		search := c.Args()
		if len(search) == 0 || strings.TrimSpace(search[0]) == "" {
			output(paths)
			return nil
		}

		rex := generateRegex(search)

		for _, path := range paths {
			if rex.MatchString(path) {
				res, err := client.Logical().Read(path)
				if err != nil || res == nil {
					continue
				}

				data := res.Data
				data["path"] = path
				output(data)

				return nil // only output the first match
			}
		}
		return nil
	}

	app.Run(os.Args)
}

func generateRegex(searchTerms []string) *regexp.Regexp {
	regex := ""
	for _, term := range searchTerms {
		regex += `.*` + term + `.*`
	}

	return regexp.MustCompile(regex)
}

func setToken(client *api.Client) error {
	token, err := ioutil.ReadFile(os.Getenv("HOME") + "/.scope/auth.token")
	if err != nil {
		return errors.New("Please provide a token file at ~/.scope/auth.token")
	}

	client.SetToken(strings.TrimSpace(string(token)))
	return nil
}

func output(data interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(data)
}
