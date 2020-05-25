package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strconv"
)

type MessageCode byte

type Peer struct {
	conn net.Conn
	Ip   net.IP
	Port uint16
}

type ResultPiece struct {
	id   int
	data []byte
}

func (p *Peer) handshake(id string, infoHash []byte) {

	var err error

	p.conn, err = net.Dial("tcp", "127.0.0.1:58560")

	if err != nil {
		log.Fatal(err)
	}

	//defer p.conn.Close()

	buffer := bytes.Buffer{}

	buffer.Write([]byte{0x13})
	buffer.Write([]byte("BitTorrent protocol"))
	buffer.Write([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	buffer.Write(infoHash)
	buffer.Write([]byte(id))

	write, err := p.conn.Write(buffer.Bytes())

	if err != nil {
		log.Fatal(write)
	}

	log.Printf("Hey I have wrote something %d", buffer.Bytes())

	readBuffer := make([]byte, 68)

	_, err = p.conn.Read(readBuffer)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Hey I have read something %d", readBuffer)

	readBuffer = make([]byte, 173)

	_, err = p.conn.Read(readBuffer)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("bitfield %d", readBuffer)

}

func (p *Peer) String() string {
	return p.Ip.String() + ":" + strconv.Itoa(int(p.Port))
}

func (p *Peer) work(pieces <-chan RequestPiece, resultPieces chan ResultPiece) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("### Recovered", r)
		}
	}()

	if p.conn == nil {
		log.Panic("peer not connected")
	}

	p.write(NewInterested())
	p.write(NewUnchoked())

	// unchocked
	p.read()

	for piece := range pieces {

		numBytesDownloaded := 0

		pieceResult := make([]byte, piece.length, piece.length)

		for numBytesDownloaded < piece.length {

			log.Printf("downloaded %d bytes for piece %d", numBytesDownloaded, piece.id)

			length := calculateLength(piece, numBytesDownloaded)

			log.Print("length calculated ", length)

			p.write(NewRequest(piece.id, numBytesDownloaded, length))

			pieceMessage := p.read()

			if pieceMessage.Code() == 7 {
				//index := binary.BigEndian.Uint32(pieceMessage.Payload()[0:4])
				offset := binary.BigEndian.Uint32(pieceMessage.Payload()[4:8])
				data := pieceMessage.Payload()[8:]

				copy(pieceResult[offset:], data)
				numBytesDownloaded += len(data)
			}

		}

		log.Printf("downloaded %d bytes for piece %d", numBytesDownloaded, piece.id)

		log.Printf("piece %d completed", piece.id)

		resultPieces <- ResultPiece{id: piece.id, data: pieceResult[:]}

	}
}

func calculateLength(piece RequestPiece, received int) int {
	left := piece.length - received

	if left < 16384 {
		return left
	}

	return 16384

}

func (p *Peer) write(message Message) {

	marshallMessage := message.Marshall()

	fullMessage := make([]byte, len(marshallMessage)+4)

	binary.BigEndian.PutUint32(fullMessage[0:4], uint32(message.PayloadLength()+1))

	copy(fullMessage[4:], marshallMessage)

	//log.Printf("sending message %d to peer %s", message.Code(), p)

	_, _ = p.conn.Write(fullMessage)

}

func (p *Peer) read() Message {

	buf := make([]byte, 4)

	_, _ = p.conn.Read(buf)
	nextMessageSize := binary.BigEndian.Uint32(buf)

	if nextMessageSize == 0 { // 0 is keep alive read again
		_, _ = p.conn.Read(buf)
		nextMessageSize = binary.BigEndian.Uint32(buf)
	}

	unParsedMessage := make([]byte, nextMessageSize, nextMessageSize)

	_, _ = p.conn.Read(unParsedMessage)

	message := Parse(unParsedMessage)

	//log.Printf("receive message %d from peer %s", message.Code(), p)

	return message
}
