package p2p

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

// MaxBlockSize is the largest number of bytes a request can ask for
const MaxBlockSize = 16384

// MaxBacklog is the number of unfulfilled requests a client can have in its pipeline
const MaxBacklog = 5

type Torrent struct {
	Peers        []Peer
	PeerID       [20]byte
	InfoHash     [20]byte
	PiecesHashes [][20]byte
	PieceLength  int
	Length       int
	Name         string
}
type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}
type pieceProgress struct {
	index      int
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}
type Peer struct {
	IP   net.IP
	Port uint16
}

func UnmarshalPeers(peerResp []byte) ([]Peer, error) {
	peers := make([]Peer, len(peerResp)/6)

	if len(peerResp)%6 != 0 {
		return nil, fmt.Errorf("Incorrect peer response size")
	}

	for i := 0; i < len(peerResp)/6; i++ {
		peers[i].IP = net.IP(peerResp[i*6 : (i*6)+4])
		peers[i].Port = binary.BigEndian.Uint16(peerResp[i*6+4 : (i*6)+6])
	}

	return peers, nil
}

func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}
