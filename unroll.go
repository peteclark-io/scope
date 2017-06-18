package main

import (
	"encoding/json"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/vault/api"
)

type list struct {
	Keys []string `json:"keys"`
}

func unroll(startingPath string, client *api.Client) ([]string, error) {
	root, found, err := listPath(startingPath, client)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errors.New(startingPath + " not found")
	}

	return listAllPaths(startingPath, root, client), nil
}

func listAllPaths(path string, root *list, client *api.Client) []string {
	results := make([]string, 0)

	for _, p := range root.Keys {
		result, found, err := listPath(path+p, client)
		if err != nil {
			log.WithField("path", p).WithError(err).Warn("Failed to parse path.")
			continue
		}

		if !found {
			results = append(results, path+p)
			continue
		}

		next := listAllPaths(path+p, result, client)
		results = append(results, next...)
	}
	return results
}

func listPath(path string, client *api.Client) (*list, bool, error) {
	res, err := client.Logical().List(path)
	if err != nil {
		log.WithError(err).Info("oh no")
		return nil, false, err
	}

	if res == nil {
		return nil, false, nil
	}

	d, _ := json.Marshal(res.Data)

	l := &list{}
	err = json.Unmarshal(d, &l)
	if err != nil {
		log.WithError(err).Warn("oh no")
		return nil, true, err
	}

	return l, true, nil
}
