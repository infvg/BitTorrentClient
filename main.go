package main

import (
	"client/torrent"
	"fmt"
)

func main() {
	//const port uint16 = 6925
	//peerID := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	//fmt.Print(os.Stat("[YURI] Kimi no Suizou o Tabetai [BD1080p HEVC FLAC][Dual Audio]  v4.torrent"))

	torrent, err := torrent.ToTorrentFile("[Yameii] Attack on Titan The Final Season - 28 [English Dub] [WEB-DL 1080p] [D3857496].mkv.torrent")

	fmt.Println(torrent.Trackers)
	if err != nil {
		return
	}
	//info2 := torrent.Open("CC_1914_08_31_TheGoodforNothing_archive.torrent")

	//fmt.Println(info2.InfoHash)

	//fmt.Print(" Line ")
	//fmt.Print(info2.RequestPeers(peerID, port))
	//fmt.Print(info2.RequestPeers(peerID, port))

	//fmt.Print(info2.RequestPeers(peerID, port))

	//torrent.TorrentFile.Announce
}

// copy all the methods he has
