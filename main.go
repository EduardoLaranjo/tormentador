package main

import (
	"log"
	"os"
)

func main() {

	currDir, err := os.Getwd()

	if err != nil {
		log.Fatal("failed to get current directory")
	}

	beeTorrent := Open(currDir + "/resources/debian-iso.torrent")

	torrent := NewTorrent(beeTorrent)

	torrent.Download()

}
