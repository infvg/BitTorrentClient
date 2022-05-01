package main

import (
	"client/torrent"
	"fmt"
)

func main() {
	const port uint16 = 6925
	peerID := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	//fmt.Print(os.Stat("[YURI] Kimi no Suizou o Tabetai [BD1080p HEVC FLAC][Dual Audio]  v4.torrent"))

	info2 := torrent.Open("[YURI] Kimi no Suizou o Tabetai [BD1080p HEVC FLAC][Dual Audio]  v4.torrent")

	fmt.Print("Line  16")

	fmt.Print(info2.RequestPeers(peerID, port))

	//torrent.TorrentFile.Announce
}

// copy all the methods he has
