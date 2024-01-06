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

/*
	// HelloFS example

import (

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
*/
import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type Record struct {
	PermData uint16
	UID      uint32
	GID      uint32
	Pointer  uint64
	Size     uint64
	Ctime    int64
	Mtime    int64
	Atime    int64
	Checksum [16]byte
	Name     [255]byte
}

func ProcessFilesAndFolders(rootDir string) error {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		var record Record

		if len(info.Name()) <= 255 {
			copy(record.Name[:], []byte(info.Name()))
		}

		record.Size = uint64(info.Size())

		record.Ctime, _ = info.Sys().(*syscall.Stat_t).Ctim.Unix()
		record.Mtime, _ = info.Sys().(*syscall.Stat_t).Mtim.Unix()
		record.Atime, _ = info.Sys().(*syscall.Stat_t).Atim.Unix()

		record.GID = info.Sys().(*syscall.Stat_t).Gid
		record.UID = info.Sys().(*syscall.Stat_t).Uid

		record.PermData = 1

		if info.Mode()&os.ModeDir != 0 {
			record.PermData = 2
		} else if info.Mode()&os.ModeSymlink != 0 {
			record.PermData = 3
		}

		if record.PermData == 1 { // Optimize for streams
			contents, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Error:", err)
				return nil
			}
			record.Checksum = md5.Sum(contents)
		}

		record.PermData = (record.PermData << 13)

		//record.PermData = (record.PermData << 11)

		record.PermData += uint16(info.Mode().Perm())

		fmt.Printf("Permdata: %016b  Time: M %016x A %016x C %016x  GID: %08x UID: %08x  Checksum: %x  Size: %d \tName: %s\n",
			record.PermData, record.Mtime, record.Atime, record.Ctime, record.GID, record.UID, record.Checksum, record.Size, record.Name)

		return nil
	})

	return err
}

func main() {
	//var record Record
	//fmt.Println(record)
	ProcessFilesAndFolders("testdir")
}
