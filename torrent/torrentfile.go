package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/url"
	"os"
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

func (b *bencodeTorrent) ToTorrentFile() (TorrentFile, error) {
	ber, err := b.Info.hash()
	if err != nil {
		return TorrentFile{}, nil
	}
	hashes, err := b.Info.splitHash()
	if err != nil {
		return TorrentFile{}, nil
	}
	t := TorrentFile{
		Announce:    b.Announce,
		InfoHash:    ber,
		PiecesHash:  hashes,
		PieceLength: b.Info.PieceLength,
		Comment:     b.Comment,
		Creation:    b.Creation,
		Length:      b.Info.Length,
		Name:        b.Info.Name,
	}
	return t, nil
}
func Open(path string) (TorrentFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer f.Close()
	b := bencodeTorrent{}
	err = bencode.Unmarshal(f, &b)
	if err != nil {
		return TorrentFile{}, err
	}
	return b.ToTorrentFile()
}
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
