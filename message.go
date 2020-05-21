package main

import (
	"encoding/binary"
)

const (
	choke MessageCode = iota
	unchoke
	interested
	notInterested
	have
	bitfield
	request
	piece
	cancel
)

func Parse(unparsedMessage []byte) Message {

	message := message{action: MessageCode(unparsedMessage[0]), payload: []byte{}}

	if message.Code() > 4 {
		message.payload = unparsedMessage[1:]
	}

	return &message

}

func NewUnchoked() Message {
	return &message{action: unchoke}
}

func NewInterested() Message {
	return &message{action: interested}
}

func NewRequest(pieceId int, offset int, length int) Message {
	payload := make([]byte, 12)

	binary.BigEndian.PutUint32(payload[0:4], uint32(pieceId))
	binary.BigEndian.PutUint32(payload[4:8], uint32(offset))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))

	return &message{
		action:  request,
		payload: payload,
	}

}

type message struct {
	action  MessageCode
	payload []byte
}

type Message interface {
	Code() MessageCode
	Payload() []byte
	PayloadLength() int
	Marshall() []byte
}

func (r *message) Code() MessageCode {
	return r.action
}

func (r *message) Payload() []byte {
	return r.payload
}

func (r *message) PayloadLength() int {
	return len(r.payload)
}

func (r *message) Marshall() []byte {
	result := make([]byte, r.PayloadLength()+1)
	result[0] = byte(r.Code())
	copy(result[1:], r.Payload())
	return result
}
