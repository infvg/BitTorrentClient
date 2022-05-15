package main

import (
	"client/torrent"
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

	fmt.Println(torrentFile.Trackers)

	if err != nil {
		return
	}

	//torrent.TrackerURL(torrentFile.Trackers[1])

}
