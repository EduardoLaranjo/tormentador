package main

import (
	"log"
	"os"
	"testing"
)

func Test_openFile(t *testing.T) {

	_, err := os.Getwd()

	if err != nil {
		log.Fatal("failed to get current directory")
	}

	torrent := BeeTorrent{}
	torrent.parse()

	//torrent := Open(currDir + "/resources/debian-iso.torrent")

}
