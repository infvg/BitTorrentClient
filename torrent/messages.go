package torrent

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
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

func (peerClient *Client) GetPiece(pieceIndex, length int, hash [20]byte) ([]byte, error) {

	const maxBlockSize = 16384
	const maxBacklog = 15

	if !peerClient.bitfield.HasPiece(pieceIndex) {
		return nil, errors.New("Peer does not have the requested piece")
	}

	peerClient.Connection.SetDeadline(time.Now().Add(time.Second * 15))
	defer peerClient.Connection.SetDeadline(time.Time{})

	var requested, received, backlog int
	pieceBuffer := make([]byte, length)
	for received < length {

		for !peerClient.isChoked && backlog < maxBacklog && requested < length {

			payload := make([]byte, 12)
			binary.BigEndian.PutUint32(payload[0:4], uint32(pieceIndex))
			binary.BigEndian.PutUint32(payload[4:8], uint32(requested))

			blockSize := maxBlockSize

			if requested+blockSize > length {
				blockSize = length - requested
			}
			binary.BigEndian.PutUint32(payload[8:12], uint32(blockSize))

			err := peerClient.SendMessage(MsgRequest, payload)
			if err != nil {
				return nil, fmt.Errorf("Creating Backlog: %s", err)
			}
			requested += blockSize
			backlog++
		}

		if peerClient.isChoked {
			err := peerClient.SendMessage(MsgUnchoke, nil)
			if err != nil {
				return nil, fmt.Errorf("Unchoke peer client: %s", err)
			}
		}

		msg, err := peerClient.RecieiveMessage()
		if err != nil {
			return nil, fmt.Errorf("receiving piece message from peer: %w", err)
		}

		if msg.id != MsgPiece {
			continue
		}
		respIndex := binary.BigEndian.Uint32(msg.payload[0:4])
		if respIndex != uint32(pieceIndex) {
			continue
		}

		start := binary.BigEndian.Uint32(msg.payload[4:8])
		blockData := msg.payload[8:]
		n := copy(pieceBuffer[start:], blockData[:])

		if n != 0 {
			received += n
			backlog--
		}
	}

	pieceHash := sha1.Sum(pieceBuffer)
	if !bytes.Equal(pieceHash[:], hash[:]) {
		return nil, fmt.Errorf("InfoHash dont match from %s", peerClient.Connection.RemoteAddr())
	}

	return pieceBuffer, nil
}
