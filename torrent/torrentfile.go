package torrent

import (
	"io"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Creation int         `bencode:"creation date"`
	Comment  string      `bencode:"comment"`
	Info     bencodeInfo `bencode:"info"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PiecesHash  [][20]byte
	PieceLength int
	Comment     string
	Creation    int
	Length      int
	Name        string
}

func Open(r io.Reader) (*bencodeTorrent, error) {
	b := bencodeTorrent{}
	err := bencode.Unmarshal(r, &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
func (b bencodeInfo) ToTorrentFile() (TorrentFile, error) {
	return TorrentFile{}, nil
}

func (t TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {

	u, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	u.RawQuery = params.Encode()
	return u.String(), nil
}
