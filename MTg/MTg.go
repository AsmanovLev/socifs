package MTg

import (
	"bytes"
	"io"
	"log"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

const (
	savechat = &tg.InputPeerSelf{}
)

func TgUpload(tgc TelegramClient, payload []byte, name string) { // byte - payload

	c := tgc.client.CreateContext()
	u := uploader.NewUploader(tgc.api)
	sender := message.NewSender(tgc.api).WithUploader(u)
	target := sender.To(savechat)

	upload, err := u.FromBytes(c, name, payload)

	if err != nil {
		panic(err)
	}

	document := message.UploadedDocument(upload, styling.Plain(``))

	document.Filename(name).ForceFile(true)

	if _, err := target.Media(c, document); err != nil {
		panic(err)
	}
}

func TgSearch(tgc TelegramClient, query string) tg.Document {
	c := tgc.client.CreateContext()
	res, err := tgc.api.MessagesSearch(c,
		&tg.MessagesSearchRequest{
			Q:      query,
			Peer:   savechat,
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

	res.Encode(buf)
	slice.Decode(buf)

	slice.Messages[0].Encode(buf)
	message.Decode(buf)

	message.Media.Encode(buf)
	media.Decode(buf)

	media.Document.Encode(buf)
	document.Decode(buf)

	/* TODO: name check
	attribute := tg.DocumentAttributeFilename{}
	document.Attributes[0].Encode(buf)
	attribute.Decode(buf)
	fmt.Println(attribute.FileName)
	*/

	return document
}

func TgDownload(tgc TelegramClient, doc tg.Document) []byte {
	c := tgc.client.CreateContext()
	d := downloader.NewDownloader()
	loc := doc.AsInputDocumentFileLocation()
	var buf bytes.Buffer
	writer := io.Writer(&buf)
	_, err := d.Download(tgc.api, loc).Stream(c, writer)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func TgInit(config ConfigTelegram) TelegramClient {
	var tgClient TelegramClient
	var err error
	tgClient.client, err = gotgproto.NewClient(
		config.AppID,
		config.AppHash,
		gotgproto.ClientType{
			Phone: config.Phone,
		},
		&gotgproto.ClientOpts{
			Session: sessionMaker.SqliteSession(config.Session),
		},
	)

	if err != nil {
		log.Fatalln("failed to start client:", err)
	}
	tgClient.api = tgClient.client.API()
	return tgClient
}

type TelegramClient struct {
	client *gotgproto.Client
	api    *tg.Client
}
type ConfigTelegram struct {
	AppID   int    `yaml:"appid"`
	AppHash string `yaml:"apphash"`
	Session string `yaml:"session"`
	Phone   string `yaml:"phone"`
}
