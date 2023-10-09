package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	. "go-http-replicator/replicator"
)

type Config struct {
	Type     string   `json:"type"`
	FetchUrl string   `json:"fetch_url"`
	Targets  []string `json:"targets"`
}

func FetchTargets(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching targets: %s\n", err)
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading targets response body: %s\n", err)
		return nil, err
	}

	var targets []string

	err = json.Unmarshal(body, &targets)
	if err != nil {
		fmt.Printf("Error parsing targets response body: %s\n", err)
		return nil, err
	}

	return targets, nil
}

func LoadConfiguration(file string) (*Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return &config, nil
}

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	config, err := LoadConfiguration(filepath.Join(path, "config.json"))
	if err != nil {
		log.Fatalln(err)
	}

	var replicator Replicator

	if config.Type == "dynamic" {
		targets, err := FetchTargets(config.FetchUrl)
		if err != nil {
			log.Fatalln(err)
		}
		replicator.SetTargets(targets)
	} else {
		replicator.SetTargets(config.Targets)
	}

	http.HandleFunc("/", replicator.Handler)
	err = http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
