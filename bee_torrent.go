package main

import (
	"bytes"
	"crypto/sha1"
	"github.com/jackpal/bencode-go"
	"log"
	"os"
)

type info struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type BeeTorrent struct {
	Info     info
	Announce string
	Comment  string
}

func Open(path string) BeeTorrent {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal("failed to get current directory")
	}

	torrentFile := BeeTorrent{}

	err = bencode.Unmarshal(file, &torrentFile)

	if err != nil {
		log.Fatal("failed to decode torrent file")
	}

	return torrentFile
}

func (b *BeeTorrent) infoHash() [20]byte {
	buffer := bytes.Buffer{}

	err := bencode.Marshal(&buffer, b.Info)

	if err != nil {
		log.Fatal(err)
	}

	return sha1.Sum(buffer.Bytes())

}
