package torrent

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zeebo/bencode"
)

type HTTPtrackerResponce struct {
	Interval   int    `bencode:"interval"`
	Peers      string `bencode:"peers"`
	Complete   int    `bencode:"complete"`
	Incomplete int    `bencode:"incomplete"`
}
type UDPstats struct {
	interval uint32
	seeders  uint32
	leecher  uint32
}

type UDPClient struct {
	conn         *net.UDPConn
	peerID       [20]byte
	infoHash     [20]byte
	port         int
	peers        []net.TCPAddr
	connectionID uint64
}
type udpMessageAction uint32

const (
	Connecting udpMessageAction = iota
	AnnounceUDP
	Error
)

var actionStrings = map[udpMessageAction]string{
	Connecting:  "connect",
	AnnounceUDP: "announce",
	Error:       "error",
}

func (m udpMessageAction) String() string {
	return actionStrings[m]
}

func TrackerURL(trackerURL string, infoHash, peerID [20]byte, port int) ([]net.TCPAddr, error) { //gets the info per tracker URL, so we gotta call
	//it for each trackers

	link, err := url.Parse(trackerURL)
	if err != nil {
		fmt.Errorf("Error in making trackerURL")
	}

	if link.Scheme == "http" || link.Scheme == "https" {
		parameters := url.Values{
			"info_hash":  []string{string(infoHash[:])},
			"peer_id":    []string{string(peerID[:])},
			"port":       []string{strconv.Itoa(int(port))},
			"uploaded":   []string{"0"},
			"downloaded": []string{"0"},
			"compact":    []string{"1"},
			"left":       []string{"0"},
		}
		link.RawQuery = parameters.Encode()

		//making the request
		getTrackerInfo := &http.Client{Timeout: 4 * time.Second}
		resp, err := getTrackerInfo.Get(link.String())

		if err != nil {
			fmt.Errorf("Error in GET Tracker")
		}

		defer resp.Body.Close()

		trackRes := HTTPtrackerResponce{}

		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("Error in Reading tracker responce")
		}

		err = bencode.DecodeBytes(raw, &trackRes)

		var peerAddresses []net.TCPAddr

		const peerSize = 6 // 4 bytes for IP, 2 for Port

		if len(trackRes.Peers)%peerSize != 0 {
			return nil, fmt.Errorf("malformed http tracker response: %w", err)
		}

		for i := 0; i < len(trackRes.Peers); i += peerSize {
			// convert port substring into byte slice to calculate via BigEndian
			portRaw := []byte(trackRes.Peers[i+4 : i+6])

			port := binary.BigEndian.Uint16(portRaw)

			peerAddresses = append(peerAddresses, net.TCPAddr{

				IP:   []byte(trackRes.Peers[i : i+4]),
				Port: int(port),
			})

		}

		return peerAddresses, nil

	} else if link.Scheme == "udp" {

		udpAddr, err := net.ResolveUDPAddr(link.Scheme, link.Host)
		if err != nil {
			return nil, fmt.Errorf("Resolving UDP address: %s", err)
		}

		conn, err := net.DialUDP("udp", nil, udpAddr)

		err = conn.SetReadBuffer(4096)
		if err != nil {
			return nil, fmt.Errorf("UDP conn read buffer: %s", err)
		}

		udpClient := UDPClient{
			conn:     conn,
			peerID:   peerID,
			infoHash: infoHash,
			port:     port,
		}

		udpClient.UDPconnection()

		return udpClient.peers, nil

	} else {

		return nil, fmt.Errorf("Wrong tracker url scheme: %s", link.Scheme)
	}

}

func (UDPclient *UDPClient) UDPconnection() error {

	const magicNumber = 0x41727101980

	transactionID := rand.Uint32()
	connectMsg := make([]byte, 16)
	binary.BigEndian.PutUint64(connectMsg[0:8], uint64(magicNumber))
	binary.BigEndian.PutUint32(connectMsg[8:12], uint32(Connecting))
	binary.BigEndian.PutUint32(connectMsg[12:16], transactionID)

	UDPclient.conn.SetDeadline(time.Now().Add(time.Second * 5))

	defer UDPclient.conn.SetDeadline(time.Time{})

	_, err := UDPclient.conn.Write(connectMsg)

	if err != nil {
		return fmt.Errorf("sending connect message to UDP: %s", err)
	}

	// Read connect response from server
	resp := make([]byte, 16)
	n, _, _, _, err := UDPclient.conn.ReadMsgUDP(resp, nil)
	if err != nil {
		return fmt.Errorf("reading connect responce: %s", err)
	}
	if n != 16 {
		return fmt.Errorf("want connect message to be 16 bytes, got %d", n)
	}

	connectResp, err := UDPclient.parseUDPResponse(transactionID, Connecting, resp)
	if err != nil {
		return fmt.Errorf("udp connect responce: %s", err)
	}

	UDPclient.connectionID = binary.BigEndian.Uint64(connectResp)

	UDPclient.announce()

	return nil

}

func (UDPclient *UDPClient) announce() error {

	UDPtrackerMessage := make([]byte, 98)
	transactionID := rand.Uint32()

	binary.BigEndian.PutUint64(UDPtrackerMessage[0:8], UDPclient.connectionID)
	binary.BigEndian.PutUint32(UDPtrackerMessage[8:12], uint32(AnnounceUDP))
	binary.BigEndian.PutUint32(UDPtrackerMessage[12:16], transactionID)
	copy(UDPtrackerMessage[16:36], UDPclient.infoHash[:])
	copy(UDPtrackerMessage[36:56], UDPclient.peerID[:])

	binary.BigEndian.PutUint64(UDPtrackerMessage[56:64], 0)
	binary.BigEndian.PutUint64(UDPtrackerMessage[64:72], 0)
	binary.BigEndian.PutUint64(UDPtrackerMessage[72:80], 0)

	binary.BigEndian.PutUint32(UDPtrackerMessage[80:84], 0)
	binary.BigEndian.PutUint32(UDPtrackerMessage[84:88], 0)

	binary.BigEndian.PutUint32(UDPtrackerMessage[88:92], rand.Uint32())

	num := -1
	binary.BigEndian.PutUint32(UDPtrackerMessage[92:96], uint32(num))
	binary.BigEndian.PutUint16(UDPtrackerMessage[96:98], uint16(UDPclient.port))

	UDPclient.conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer UDPclient.conn.SetDeadline(time.Time{})

	_, err := UDPclient.conn.Write(UDPtrackerMessage)
	if err != nil {
		return fmt.Errorf("Writing announce message: %s", err)
	}

	responce := make([]byte, 4096)
	n, err := UDPclient.conn.Read(responce)
	if err != nil {
		return fmt.Errorf("Reading announce response: %s", err)
	}

	responce = responce[:n]
	announceResp, err := UDPclient.parseUDPResponse(transactionID, AnnounceUDP, responce)
	if err != nil {
		return fmt.Errorf("UDP announce responce: %s", err)
	}

	interval := binary.BigEndian.Uint32(announceResp[0:4])
	leecher := binary.BigEndian.Uint32(announceResp[4:8])
	seeders := binary.BigEndian.Uint32(announceResp[8:12])

	stats := UDPstats{
		interval: interval,
		seeders:  seeders,
		leecher:  leecher,
	}
	fmt.Println("Tracker stats: Interval:", stats.interval, "\nSeeders:", stats.seeders, "\nLeechers:", stats.leecher) // printing stats

	var peers []net.TCPAddr
	for i := 12; i < len(announceResp); i += 6 {
		peers = append(peers, net.TCPAddr{
			IP:   announceResp[i : i+4],
			Port: int(binary.BigEndian.Uint16(announceResp[i+4 : i+6])),
		})
	}

	if len(peers) == 0 {
		return fmt.Errorf("no peers found")
	}
	UDPclient.peers = peers

	return nil
}

func (u *UDPClient) parseUDPResponse(wantTransactionID uint32, wantAction udpMessageAction, resp []byte) ([]byte, error) {
	if len(resp) < 8 {
		return nil, fmt.Errorf("response is <8 characters, got %d", len(resp))
	}

	respTransactionID := binary.BigEndian.Uint32(resp[4:8])
	if respTransactionID != wantTransactionID {
		return nil, fmt.Errorf("transactionIDs do not match, want %d, got %d", wantTransactionID, respTransactionID)
	}

	action := binary.BigEndian.Uint32(resp[0:4])
	if udpMessageAction(action) == Error {
		// return an error that includes the message
		errorText := string(resp[8:])
		return nil, fmt.Errorf("error response: %s", errorText)
	}
	if udpMessageAction(action) != wantAction {
		return nil, fmt.Errorf("want %s action, got %s", wantAction, udpMessageAction(action))
	}

	return resp[8:], nil
}
