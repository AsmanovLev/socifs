package main

/*
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
*/

import (
	"fmt"
	"github.com/winfsp/cgofuse/fuse"
	"os"
)

const (
	filename = "hello"
	contents = "hello, world\n"
)

type Hellofs struct {
	fuse.FileSystemBase
}

func (self *Hellofs) Open(path string, flags int) (errc int, fh uint64) {
	switch path {
	case "/" + filename:
		return 0, 0
	default:
		return -fuse.ENOENT, ^uint64(0)
	}
}

func (self *Hellofs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	fmt.Println("someone read attributes")
	switch path {
	case "/":
		stat.Mode = fuse.S_IFDIR | 0555
		return 0
	case "/" + filename:
		stat.Mode = fuse.S_IFREG | 0444
		stat.Size = int64(len(contents))
		return 0
	default:
		return -fuse.ENOENT
	}
}

func (self *Hellofs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	endofst := ofst + int64(len(buff))
	if endofst > int64(len(contents)) {
		endofst = int64(len(contents))
	}
	if endofst < ofst {
		return 0
	}
	n = copy(buff, contents[ofst:endofst])
	return
}

func (self *Hellofs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	fill(".", nil, 0)
	fill("..", nil, 0)
	fill(filename, nil, 0)
	return 0
}

func main() {
	hellofs := &Hellofs{}
	host := fuse.NewFileSystemHost(hellofs)
	host.Mount("", os.Args[1:])
}
