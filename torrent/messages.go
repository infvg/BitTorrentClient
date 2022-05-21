package torrent

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

type Message struct {
	id      messageID
	payload []byte
}
type bitfield []byte

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
	MsgKeepAlive     messageID = 254
	MsgPortNum       messageID = 20
)

var messageIDMap = map[messageID]string{
	MsgChoke:         "choke",
	MsgUnchoke:       "unchoke",
	MsgInterested:    "interested",
	MsgNotInterested: "not interested",
	MsgHave:          "have",
	MsgBitfield:      "bitfield",
	MsgRequest:       "request",
	MsgPiece:         "piece",
	MsgCancel:        "cancel",
	MsgKeepAlive:     "keep alive",
	MsgPortNum:       "port",
}

func (msg messageID) String() string {
	return messageIDMap[msg]
}

func (peerClient *Client) SendMessage(id messageID, payload []byte) error {

	length := uint32(len(payload) + 1)
	message := make([]byte, length+4)
	binary.BigEndian.PutUint32(message[0:4], length)
	message[4] = byte(id)
	copy(message[5:], payload)

	_, err := peerClient.Connection.Write(message)
	if err != nil {
		return fmt.Errorf("Sending message: %s", err)
	}

	return nil
}

func (peerClient *Client) RecieiveMessage() (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(peerClient.Connection, lengthBuf)

	if err != nil {
		return nil, err
	}

	messageLenght := binary.BigEndian.Uint32(lengthBuf)

	if messageLenght == 0 {
		message := Message{
			id:      MsgKeepAlive,
			payload: nil,
		}

		return &message, nil
	}

	messageBuf := make([]byte, messageLenght)
	_, err = io.ReadFull(peerClient.Connection, messageBuf)
	if err != nil {
		return nil, err
	}

	message := Message{
		id:      messageID(messageBuf[0]),
		payload: messageBuf[1:],
	}

	switch message.id {

	case MsgChoke:
		peerClient.isChoked = true

	case MsgUnchoke:
		peerClient.isChoked = false

	case MsgHave:
		index := binary.BigEndian.Uint32(message.payload)
		peerClient.bitfield.SetPiece(int(index))

	case MsgBitfield:
		peerClient.bitfield = bitfield(message.payload)

	case MsgPortNum:
		peerClient.dhtPort = int(binary.BigEndian.Uint16(message.payload))
	}

	return &message, nil
}

func (bf bitfield) HasPiece(index int) bool {
	byteIndex := index / 8

	offset := index % 8
	return bf[byteIndex]>>(7-offset)&1 != 0

}

func (bf bitfield) SetPiece(index int) {

	byteIndex := index / 8

	offset := index % 8
	mask := 1 << (7 - offset)
	bf[byteIndex] |= byte(mask)

}
