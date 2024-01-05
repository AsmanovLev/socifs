package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"

	//"github.com/gotd/td/bin"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

// generate key with dd if=/dev/random of=aes256.key bs=256 count=1

func upload(client *gotgproto.Client, api *tg.Client, payload []byte, name string) { // byte - payload

	c := client.CreateContext()
	u := uploader.NewUploader(api)
	sender := message.NewSender(api).WithUploader(u)
	target := sender.To(&tg.InputPeerSelf{})

	upload, err := u.FromBytes(c, "test.txt", payload)
	log.Println("Uploading file")
	if err != nil {
		fmt.Errorf("upload %q: %w", "binary", err)
		return
	}

	document := message.UploadedDocument(upload, styling.Plain(`Upload: From bot`))

	document.Filename(name).ForceFile(true)

	log.Println("Sending file")
	if _, err := target.Media(c, document); err != nil {
		fmt.Errorf("send: %w", err)
		return
	}
}

func search(client *gotgproto.Client, api *tg.Client, query string) tg.Document {
	c := client.CreateContext()
	res, err := api.MessagesSearch(c,
		&tg.MessagesSearchRequest{
			Q:      query,
			Peer:   &tg.InputPeerSelf{},
			Filter: &tg.InputMessagesFilterDocument{},
			Limit:  1,
		},
	)

	if err != nil {
		panic(err)
	}

	buf := new(bin.Buffer)
	slice := tg.MessagesMessagesSlice{}
	message := tg.Message{}
	media := tg.MessageMediaDocument{}
	document := tg.Document{}
	attribute := tg.DocumentAttributeFilename{}

	res.Encode(buf)
	slice.Decode(buf)

	slice.Messages[0].Encode(buf)
	message.Decode(buf)

	message.Media.Encode(buf)
	media.Decode(buf)

	media.Document.Encode(buf)
	document.Decode(buf)

	//fmt.Println(document.AccessHash, document.ID)

	document.Attributes[0].Encode(buf)
	attribute.Decode(buf)

	fmt.Println(attribute.FileName)

	return document
}

func download(client *gotgproto.Client, api *tg.Client, doc tg.Document) {
	c := client.CreateContext()
	d := downloader.NewDownloader()
	loc := doc.AsInputDocumentFileLocation()

	_, err := d.Download(api, loc).Stream(c, os.Stdout)
	if err != nil {
		panic(err)
	}
	fmt.Println("\ndownload done")
}

type ConfigTelegram struct {
	AppID   int    `yaml:"appid"`
	AppHash string `yaml:"apphash"`
	Session string `yaml:"session"`
	Phone   string `yaml:"phone"`
}
type Config struct {
	Telegram ConfigTelegram `yaml:"telegram"`
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

	client, err := gotgproto.NewClient(
		config.Telegram.AppID,
		config.Telegram.AppHash,
		gotgproto.ClientType{
			Phone: config.Telegram.Phone,
		},
		&gotgproto.ClientOpts{
			Session: sessionMaker.SqliteSession(config.Telegram.Session),
		},
	)

	if err != nil {
		log.Fatalln("failed to start client:", err)
	}

	api := client.API()

	//log.Println("Trying to upload")

	//upload(client, api, []byte("Hello my file"))

	document := search(client, api, "journal")
	download(client, api, document) //"journal")

	//var me tg.User = Me()

	//fmt.Printf("%q", me)
	//client.Idle()
	//c, err := client.Start(&gotgproto.ClientOpts{})

}
