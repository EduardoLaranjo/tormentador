package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"github.com/jackpal/bencode-go"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type TrackerResponse struct {
	Interval int
	Peers    string
}

type Torrent struct {
	infoHash    [20]byte
	tracker     string
	id          string
	left        string
	port        int
	downloaded  int
	pieces      []SHA1
	pieceLength int
}

type RequestPiece struct {
	id     int
	length int
	hash   SHA1
}

type SHA1 []byte

const DefaultPort = 6883

func NewTorrent(beeTorrent BeeTorrent) Torrent {

	id := genId()
	infoHash := beeTorrent.infoHash()
	pieces := parsePieces(beeTorrent.Info.Pieces)

	log.Printf("my id is %s\n", id)
	log.Printf("this torrent has %d pieces\n", len(pieces))

	return Torrent{
		id:          id,
		infoHash:    infoHash,
		tracker:     beeTorrent.Announce,
		port:        DefaultPort,
		left:        "100",
		downloaded:  0,
		pieces:      pieces,
		pieceLength: beeTorrent.Info.PieceLength,
	}

}

func (t *Torrent) Download() {

	requestChannel := make(chan RequestPiece, len(t.pieces))
	resultPieces := make(chan ResultPiece)

	peers := []Peer{{Port: 58560, Ip: []byte("127.0.0.1")}}

	for _, peer := range peers {
		peer.handshake(t.id, t.infoHash[:])
		time.Sleep(time.Second)
		go peer.work(requestChannel, resultPieces)
	}

	for index, hash := range t.pieces {
		requestChannel <- RequestPiece{index, t.pieceLength, hash}
	}

	//file, err := os.OpenFile("file.iso", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//
	//buffer := bufio.NewWriter(file)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	for i := 0; i < len(t.pieces); i++ {

		piece := <-resultPieces
		resultHash := sha1.Sum(piece.data)
		isEqual := bytes.Equal(resultHash[:], t.pieces[piece.id])

		if !isEqual {
			log.Printf("integrity failed for piece %d\n", piece.id)
		}
		file, _ := os.OpenFile("debian.iso", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		log.Print("order of results ", piece.id)
		_, _ = file.Write(piece.data)
		//os.Exit(-1)
		//}
	}

	close(requestChannel)
}

func parsePieces(oldPieces string) []SHA1 {
	const PieceSize = 20

	size := len(oldPieces) / PieceSize

	var newPieces []SHA1

	for i := 0; i < size; i++ {
		offset := PieceSize * i
		newPieces = append(newPieces, []byte(oldPieces[offset:offset+PieceSize]))
	}

	return newPieces
}

func genId() string {
	builder := strings.Builder{}
	builder.WriteString(strconv.FormatInt(rand.Int63(), 9))
	return builder.String()
}

func getPeers(torrent Torrent) []Peer {
	req, err := url.Parse(torrent.tracker)

	if err != nil {
		log.Fatal(err)
	}

	query := url.Values{
		"info_hash":  []string{string(torrent.infoHash[:])},
		"peer_id":    []string{torrent.id},
		"port":       []string{string(torrent.port)},
		"uploaded":   []string{string(0)},
		"downloaded": []string{string(torrent.downloaded)},
		"left":       []string{torrent.left},
		"compact":    []string{string(1)},
	}

	req.RawQuery = query.Encode()

	get, err := http.Get(req.String())

	if err != nil {
		log.Fatal(err)
	}

	defer get.Body.Close()

	return parseResponse(get.Body)
}

func parseResponse(req io.Reader) []Peer {
	trackerResponse := TrackerResponse{}

	_ = bencode.Unmarshal(req, &trackerResponse)

	bytes := []byte(trackerResponse.Peers)

	length := len(bytes) / 6

	var peers []Peer

	for i := 0; i < length; i = i + 6 {
		ip := net.IP(bytes[i : i+4])
		port := binary.BigEndian.Uint16(bytes[i+4 : i+6])
		peers = append(peers, Peer{Ip: ip, Port: port})
	}

	return peers
}
