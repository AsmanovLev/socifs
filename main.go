package main

import (
	"log"
	"os"

	MTg "github.com/AsmanovLev/socifs/MTg"
	"gopkg.in/yaml.v3"
)

// generate key with dd if=/dev/random of=aes256.key bs=256 count=1

type Config struct {
	Telegram MTg.ConfigTelegram `yaml:"telegram"`
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
	tgc := MTg.TgInit(config.Telegram)

	log.Println("Uploading")
	MTg.TgUpload(tgc, []byte("Hello my file"), "journal")

	log.Println("Searching")
	document := MTg.TgSearch(tgc, "journal")

	log.Println("Downloading")
	downloaded := MTg.TgDownload(tgc, document)

	log.Println("Downloaded following string:", string(downloaded))
}
