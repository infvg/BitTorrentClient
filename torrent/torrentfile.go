package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"net/http"
	"time"

	"github.com/jackpal/bencode-go"
)

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}
type bencodeInfo struct {
	Pieces        string   `bencode:"pieces"`
	PieceLength   int      `bencode:"piece length"`
	Length        int      `bencode:"length"`
	Name          string   `bencode:"name"`
	AnnounceList2 []string `bencode:"announce-list"` // should return a list of all trackers

}

type bencodeTorrent struct {
	Info         bencodeInfo `bencode:"info"`
	Announce     string      `bencode:"announce"`
	AnnounceList []string    `bencode:"announce-list"` // should return a list of all trackers
	Creation     int         `bencode:"creation date"`

	CreatedBy string `bencode:"created by"`
	encoding  string `bencode:"encoding"`

	Comment string `bencode:"comment"`
}

type TorrentFile struct { //Torrent file format
	Announce     string
	AnnounceList [][]string
	CreatedBy    string
	encoding     string

	InfoHash    [20]byte
	PiecesHash  [][20]byte
	PieceLength int
	Comment     string
	Creation    int
	Length      int
	Name        string
} //The number of pieces is total length / piece size

func (b *bencodeTorrent) ToTorrentFile() TorrentFile {
	fmt.Println(b.Announce)
	fmt.Println(b.AnnounceList)

	ber, err := b.Info.hash()
	if err != nil {
		return TorrentFile{}
	}

	hashes, err := b.Info.splitHash()

	if err != nil {
		return TorrentFile{}
	}
	torrentFileInfo := TorrentFile{
		Announce:     b.Announce,
		AnnounceList: b.ToTorrentFile().AnnounceList,
		CreatedBy:    b.CreatedBy,
		encoding:     b.encoding,

		InfoHash:    ber,
		PiecesHash:  hashes,
		PieceLength: b.Info.PieceLength,
		Comment:     b.Comment,
		Creation:    b.Creation,
		Length:      b.Info.Length,
		Name:        b.Info.Name,
	}
	fmt.Println("LINE58")

	fmt.Println(torrentFileInfo.Announce)

	return torrentFileInfo
}

/*



















 */

func Open(path string) TorrentFile {
	filePath, err := os.Open(path)

	if err != nil {
		return TorrentFile{}
	}
	//readFIlesss, err := os.ReadFile(path)

	//fmt.Print("LINE65")

	//fmt.Print(string(readFIlesss))

	defer filePath.Close()

	decodedTorrent := bencodeTorrent{}
	err = bencode.Unmarshal(filePath, &decodedTorrent)

	if err != nil {
		fmt.Print("FileNotFound")
		return TorrentFile{}
	}

	//fmt.Print(decodedTorrent.ToTorrentFile())
	//fmt.Print(decodedTorrent.Announce)

	return decodedTorrent.ToTorrentFile()
}

func (torrentFileInfo TorrentFile) BuildTrackerURL(peerID [20]byte, port uint16) (string, error) {

	//fmt.Println(torrentFileInfo.InfoHash[:])

	torrentURL, err := url.Parse(torrentFileInfo.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(torrentFileInfo.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(torrentFileInfo.Length)},
	}
	torrentURL.RawQuery = params.Encode()

	return torrentURL.String(), nil
}

//From peers.go File
func (torrentFileInfo *TorrentFile) RequestPeers(peerID [20]byte, port uint16) ([]Peer, error) {

	fmt.Println(torrentFileInfo.Announce)
	fmt.Println(torrentFileInfo.Announce)

	url, err := torrentFileInfo.BuildTrackerURL(peerID, port)
	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: 15 * time.Second}

	resp, err := c.Get(url)

	if err != nil {

		return nil, err
	}

	defer resp.Body.Close()

	trackerResp := bencodeTrackerResp{}
	err = bencode.Unmarshal(resp.Body, &trackerResp)
	fmt.Print(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Print("Attac")

	return Unmarshal([]byte(trackerResp.Peers))
}

/*






// make it loop through all the trackers










*/

func (b *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *b)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (b *bencodeInfo) splitHash() ([][20]byte, error) {
	buf := []byte(b.Pieces)
	if len(buf)%20 != 0 {
		return nil, fmt.Errorf("Malformed pieces of length %d", len(buf))
	}
	hash := make([][20]byte, len(buf)/20)
	for i := 0; i < len(buf)/20; i++ {
		copy(hash[i][:], buf[i*20:(i+1)*20])
	}
	return hash, nil
}
