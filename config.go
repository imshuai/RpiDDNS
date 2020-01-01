package main

import "os"

import "io/ioutil"

import "encoding/json"

type config struct {
	Token             string
	Domain            string
	DomainID          string
	SubDomain         string
	SubDomainRecordID string
	Email             string
}

func (c *config) init() error {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	if path == "" {
		path = "settings.json"
	}
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	byts, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byts, c)
	if err != nil {
		return err
	}
	return nil
}
