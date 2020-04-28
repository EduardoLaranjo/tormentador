package main

import (
	"github.com/jackpal/bencode-go"
	"log"
	"net/http"
	"net/url"
)

type TrackerResponse struct {
	Interval int
	Peers    string
}

type Torrent struct {
	infoHash   [20]byte
	tracker    string
	myId       string
	left       string
	port       int
	downloaded int
	pieces     []byte
}

const DefaultPort = 6883
const GBitId = "G5577006791947779410"

func NewTorrent(beeTorrent BeeTorrent) Torrent {

	infoHash := beeTorrent.infoHash()

	return Torrent{
		myId:       GBitId,
		infoHash:   infoHash,
		tracker:    beeTorrent.Announce,
		port:       DefaultPort,
		left:       "100",
		downloaded: 0,
	}

}

func (t *Torrent) query() {

	req, err := url.Parse(t.tracker)

	if err != nil {
		log.Fatal(err)
	}

	query := url.Values{
		"info_hash":  []string{string(t.infoHash[:])},
		"peer_id":    []string{t.myId},
		"port":       []string{string(t.port)},
		"uploaded":   []string{string(0)},
		"downloaded": []string{string(t.downloaded)},
		"left":       []string{t.left},
	}

	req.RawQuery = query.Encode()

	get, err := http.Get(req.String())

	if err != nil {
		log.Fatal(err)
	}

	defer get.Body.Close()

	trackerResponse := TrackerResponse{}

	err = bencode.Unmarshal(get.Body, &trackerResponse)

	log.Println(trackerResponse)
	//log.Println([]byte(trackerResponse.Peers))

}
