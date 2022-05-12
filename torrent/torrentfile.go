package torrent

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zeebo/bencode"
)

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

type bencodeTorrent struct {
	Info         bencode.RawMessage `bencode:"info"`
	Announce     string             `bencode:"announce"`
	AnnounceList [][]string         `bencode:"announce-list"` // should return a list of all trackers
	Creation     int                `bencode:"creation date"`
	CreatedBy    string             `bencode:"created by"`
	encoding     string             `bencode:"encoding"`
	Comment      string             `bencode:"comment"`
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	Files       []struct {
		Length   int      `bencode:"length"` // length of this file
		Path     []string `bencode:"path"`   // list of subdirectories, last element is file name
		SHA1Hash string   `bencode:"sha1"`   // to validate this file
		MD5Hash  string   `bencode:"md5"`    // to validate this file
	} `bencode:"files"`
}

type File struct {
	Length   int    // length in bytes
	FullPath string // download path
	SHA1Hash string // for final validation
	MD5Hash  string // for final validation
}

type TorrentFile struct {
	Trackers    []string
	InfoHash    [20]byte
	PiecesHash  [][20]byte // SHA-1 hashes of each file piece
	PieceLength int        // number of bytes per piece
	Files       []File     // in the 1 file case, this will only have one element
	TotalLength int        // calculated as the sum of all files
	Name        string     // human readable display name (.torrent filename or magnet link dn)
} //The number of pieces is total length / piece size

/*



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


func New(source string) (TorrentFile, error) {
	return ToTorrentFile(source)

	//return TorrentFile{}, fmt.Errorf("invalid source (torrent file and magnet links supported)")
}
*/
func ToTorrentFile(path string) (TorrentFile, error) {
	path = os.ExpandEnv(path) //incase it is invalid
	filePath, err := os.Open(path)

	if err != nil {
		return TorrentFile{}, err
	}

	var benTorrend bencodeTorrent
	err = bencode.NewDecoder(filePath).Decode(&benTorrend)
	if err != nil {
		return TorrentFile{}, fmt.Errorf("unmarshalling file: %w", err)
	}

	var trackerURLs []string
	for _, list := range benTorrend.AnnounceList {
		trackerURLs = append(trackerURLs, list...)
	}
	// BEP0012, only use `announce` if `announce-list` is not present
	if len(trackerURLs) == 0 {
		trackerURLs = append(trackerURLs, benTorrend.Announce)
	}
	torFile := TorrentFile{
		Trackers: trackerURLs,
		Name:     path,
	}

	err = torFile.AppendMetadata(benTorrend.Info)
	if err != nil {
		return TorrentFile{}, fmt.Errorf("parsing metadata: %w", err)
	}

	return torFile, nil
}

func (torFile *TorrentFile) AppendMetadata(metadata []byte) error {
	var info bencodeInfo
	err := bencode.DecodeBytes(metadata, &info)
	if err != nil {
		return fmt.Errorf("unmarshalling info dict: %w", err)
	}

	// SHA-1 hash the entire info dictionary to get the info_hash
	torFile.InfoHash = sha1.Sum(metadata)

	// split the Pieces blob into the 20-byte SHA-1 hashes for comparison later
	const hashLen = 20 // length of a SHA-1 hash
	if len(info.Pieces)%hashLen != 0 {
		return errors.New("invalid length for info pieces")
	}
	torFile.PiecesHash = make([][20]byte, len(info.Pieces)/hashLen)
	for i := 0; i < len(torFile.PiecesHash); i++ {
		piece := info.Pieces[i*hashLen : (i+1)*hashLen]
		copy(torFile.PiecesHash[i][:], piece)
	}

	torFile.PieceLength = info.PieceLength

	// either Length OR Files field must be present (but not both)
	if info.Length == 0 && len(info.Files) == 0 {
		return fmt.Errorf("invalid torrent file info dict: no length OR files")
	}

	if info.Length != 0 {
		torFile.Files = append(torFile.Files, File{
			Length:   info.Length,
			FullPath: info.Name,
		})
		torFile.TotalLength = info.Length
	} else {
		for _, f := range info.Files {
			subPaths := append([]string{info.Name}, f.Path...)
			torFile.Files = append(torFile.Files, File{
				Length:   f.Length,
				FullPath: filepath.Join(subPaths...),
				SHA1Hash: f.SHA1Hash,
				MD5Hash:  f.MD5Hash,
			})
			torFile.TotalLength += f.Length
		}
	}

	return nil
}

/*
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
*/
