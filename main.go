package main

import (
	"log"
	"os"

	Mtg "/MTg"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram Mtg.ConfigTelegram `yaml:"telegram"`
}

func main() {
	log.Println("Reading config from config.yml...")
	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	log.Println("Connecting")
	tgc := Mtg.tgInit(config.Telegram)

	log.Println("Uploading")
	Mtg.tgUpload(tgc, []byte("Hello my file"), "journal")

	log.Println("Searching")
	document := Mtg.tgSearch(tgc, "journal")

	log.Println("Downloading")
	downloaded := Mtg.tgDownload(tgc, document)

	log.Println("Downloaded following string:", string(downloaded))
}
