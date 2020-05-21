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

	if p.conn == nil {
		log.Panic("peer not connected")
	}

	p.write(NewInterested())
	p.write(NewUnchoked())

	for piece := range pieces {

		pieceData := make([]byte, piece.length)

		received := 0

		for received < piece.length {

			var length int

			if (piece.length - received) > 16384 {
				length = 16_384
			} else {
				length = piece.length - received
			}

			log.Printf("request for piece %d with length %d\n", piece.id, length)

			p.write(NewRequest(piece.id, received, length))

			pieceMessage := p.read()

			if pieceMessage.Code() == 7 {
				index := binary.BigEndian.Uint32(pieceMessage.Payload()[0:4])
				begin := binary.BigEndian.Uint32(pieceMessage.Payload()[4:8])
				data := pieceMessage.Payload()[8:]

				log.Printf("got piece from %d with offset %d\n", index, begin)

				if begin >= uint32(received) {
					copy(pieceData[begin:], data)
					received = received + pieceMessage.PayloadLength()
				}
			}

		}

		log.Printf("piece %d completed", piece.id)

		resultPieces <- ResultPiece{id: piece.id, data: pieceData}

	}
}

func (p *Peer) write(message Message) {

	marshallMessage := message.Marshall()

	fullMessage := make([]byte, len(marshallMessage)+4)

	binary.BigEndian.PutUint32(fullMessage[0:4], uint32(message.PayloadLength()+1))

	copy(fullMessage[4:], marshallMessage)

	log.Printf("sending message %d to peer %s", message.Code(), p)

	_, _ = p.conn.Write(fullMessage)

}

func (p *Peer) read() Message {

	messageLength := make([]byte, 4)

	for binary.BigEndian.Uint32(messageLength) == 0 { // keep alives
		_, _ = p.conn.Read(messageLength)
	}

	unParsedMessage := make([]byte, binary.BigEndian.Uint32(messageLength), binary.BigEndian.Uint32(messageLength))

	_, _ = p.conn.Read(unParsedMessage)

	message := Parse(unParsedMessage)

	log.Printf("receive message %d from peer %s", message.Code(), p)

	return message
}
