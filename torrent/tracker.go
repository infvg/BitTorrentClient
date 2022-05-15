package torrent

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zeebo/bencode"
)

type TrackerResponce struct {
	Interval   int    `bencode:"interval"`
	Peers      string `bencode:"peers"`
	Complete   int    `bencode:"complete"`
	Incomplete int    `bencode:"incomplete"`
}

func TrackerURL(trackerURL string, infoHash, peerID [20]byte, port int) ([]net.TCPAddr, error) { //gets the info per tracker URL, so we gotta call
	//it for each trackers
	link, err := url.Parse(trackerURL)
	if err != nil {
		fmt.Errorf("Error in making trackerURL")
	}

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

	trackRes := TrackerResponce{}

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

}
