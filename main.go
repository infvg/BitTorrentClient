package main

import (
	"client/torrent"

	"crypto/rand"
	"fmt"
	"net"
	"sync"
)

/*
TO DO:
- Actually Downloading the piece
- After the interval has passed we should recheck the tracker
-Concurrency so it downloads from all trackers/Peers at the same time and with the Interval value given rechecks at those intervals
*/

func main() {

	Download("[Yameii] Attack on Titan The Final Season - 28 [English Dub] [WEB-DL 1080p] [D3857496].mkv.torrent") // gets all the peerInfo next we need to request the pieces

}

type torrentDownload struct {
	torrent     torrent.TorrentFile
	peerClients []*torrent.Client
}
type pieceJob struct {
	index  int
	length int
	hash   [20]byte
}

type pieceResult struct {
	Index     int
	FilePiece []byte
}

func Download(path string) (*torrentDownload, error) {

	var peerID [20]byte
	rand.Read(peerID[:])
	portNum := 6851
	var wg sync.WaitGroup
	var peerAddress []net.TCPAddr
	var peerClients []*torrent.Client
	var mutex sync.Mutex

	torrentFile, err := torrent.ToTorrentFile(path)

	if err != nil {
		return nil, fmt.Errorf("Error in opening torrent file: %s", err)
	}

	wg.Add(len(torrentFile.Trackers))
	for _, trackers := range torrentFile.Trackers {

		go getTrackerPeers(trackers, torrentFile.InfoHash, peerID, portNum, &wg, &peerAddress, &mutex)
	}

	wg.Wait()

	fmt.Println("Number of peers address:", len(peerAddress))
	wg.Add(len(peerAddress))

	for _, address := range peerAddress {
		address := address

		go func() {
			defer wg.Done()
			client, err := torrent.ConnetingToClient(torrentFile.InfoHash, peerID, address)
			if err != nil {
				fmt.Printf("error connecting to peer at %s: %s\n", address.String(), err.Error())
				return
			}

			mutex.Lock()
			peerClients = append(peerClients, client)
			mutex.Unlock()
		}()
	}
	wg.Wait()

	fmt.Println("Number of peers:", len(peerClients))
	fmt.Println("\nPeers:\n", peerClients)

	fmt.Println(torrentFile.Files)

	return &torrentDownload{
		torrent:     torrentFile,
		peerClients: peerClients,
	}, nil

}

func getTrackerPeers(trackerURL string, infoHash, peerID [20]byte, port int, wg *sync.WaitGroup, peerAddress *[]net.TCPAddr, mutex *sync.Mutex) {

	address, err := torrent.TrackerURL(trackerURL, infoHash, peerID, port)

	if err != nil {
		fmt.Errorf("Error retrieving peers from tracker: %s, %s\n", trackerURL, err)
	}

	mutex.Lock()
	*peerAddress = append(*peerAddress, address...)
	mutex.Unlock()
	wg.Done()

}
