package p2p

import (
	"encoding/binary"
	"fmt"
	"net"
)

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

func connect(peers []Peer) error {
}
