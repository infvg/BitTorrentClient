package torrent

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zeebo/bencode"
)

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
		Length  int      `bencode:"length"` // length of this file in bytes
		Path    []string `bencode:"path"`   // list of all subdirectories
		MD5Hash string   `bencode:"md5"`    // to validate this file
	} `bencode:"files"`
}

type File struct {
	Length  int
	Path    string
	MD5Hash string
}

type TorrentFile struct {
	Trackers    []string
	InfoHash    [20]byte
	PiecesHash  [][20]byte // sha-1 hashes of file pieces
	PieceLength int        // number of bytes per piece
	Files       []File     // all the files in a string list
	TotalLength int        // calculated as the sum of all files
	Name        string
} //The number of pieces is total length / piece size

func ToTorrentFile(path string) (TorrentFile, error) {

	filePath, err := os.Open(path)

	var benTorrend bencodeTorrent

	err = bencode.NewDecoder(filePath).Decode(&benTorrend)

	if err != nil {
		return TorrentFile{}, err
	}

	var trackerURLs []string

	for _, trackerList := range benTorrend.AnnounceList {
		trackerURLs = append(trackerURLs, trackerList...)
	}
	if len(trackerURLs) == 0 {
		trackerURLs = append(trackerURLs, benTorrend.Announce)
	}

	torrentMetadata := TorrentFile{

		Trackers: trackerURLs,
		Name:     path,
	}

	torrentMetadata.AddingInfoMeta(benTorrend.Info)

	return torrentMetadata, nil
}

func (torFile *TorrentFile) AddingInfoMeta(metadata []byte) error {

	var hashLen = 20 // length of a SHA-1 hash
	var info bencodeInfo
	bencode.DecodeBytes(metadata, &info)

	torFile.InfoHash = sha1.Sum(metadata)

	if len(info.Pieces)%hashLen != 0 {
		return fmt.Errorf("Error in Info lenght")
	}

	torFile.PiecesHash = make([][20]byte, len(info.Pieces)/hashLen) //slicing the pieces into constant lenghts

	for i := 0; i < len(torFile.PiecesHash); i++ {
		piece := info.Pieces[i*hashLen : (i+1)*hashLen]
		copy(torFile.PiecesHash[i][:], piece)
	}

	torFile.PieceLength = info.PieceLength

	if info.Length != 0 { // if there is one file
		torFile.Files = append(torFile.Files, File{
			Length: info.Length,
			Path:   info.Name,
		})
		torFile.TotalLength = info.Length
	} else { // if the torrent contains multiple files
		for _, f := range info.Files {
			subPaths := append([]string{info.Name}, f.Path...)
			torFile.Files = append(torFile.Files, File{
				Length:  f.Length,
				Path:    filepath.Join(subPaths...),
				MD5Hash: f.MD5Hash,
			})
			torFile.TotalLength += f.Length
		}
	}

	return nil
}
