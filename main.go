package main

import (
	"client/torrent"
	"crypto/rand"
	"fmt"
)

/*
TO DO:
-Make the tracker url and get all peers
-Connect to peers by handshake first and check if choke/unchoke
-Concurrency so it downloads from all trackers/Peers at the same time and with the Interval value given rechecks at those intervals
*/

func main() {

	torrentFile, err := torrent.ToTorrentFile("[Yameii] Attack on Titan The Final Season - 28 [English Dub] [WEB-DL 1080p] [D3857496].mkv.torrent")
	var peerID [20]byte
	rand.Read(peerID[:])
	port := 6851

	fmt.Println(torrentFile.Trackers)

	if err != nil {
		return
	}

	for count := 1; count < len(torrentFile.Trackers); count++ {
		fmt.Println("Tracker number: ", count)
		fmt.Print(torrent.TrackerURL(torrentFile.Trackers[count], torrentFile.InfoHash, peerID, port))

	}

}
