package main

import (
	"client/torrent"
	"os"
	"path/filepath"
	"time"

	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"sync"
)

/*
TO DO:
- After the interval has passed we should recheck the tracker
- Concurrency so it downloads from all trackers/Peers at the same time and with the Interval value given rechecks at those intervals
- Gives the Download speed
- UI
- Writing data to file has pieces are retrieved
*/

func main() {

	ready, err := DownloadInfo("[Yameii] Attack on Titan The Final Season - 28 [English Dub] [WEB-DL 1080p] [D3857496].mkv.torrent") // gets all the peerInfo next we need to request the pieces

	if err != nil {
		fmt.Println(err)
	}
	ready.DownloadFiles()

}

type torrentDownloadInfo struct {
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

func DownloadInfo(path string) (*torrentDownloadInfo, error) {

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

	for _, file := range torrentFile.Files {
		fmt.Printf("File Name: %s\nSize: %.2f GB", file.Path, float64(file.Length)/(1024*1024*1024))

	}
	time.Sleep(4 * time.Second)

	return &torrentDownloadInfo{
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

func (files *torrentDownloadInfo) DownloadFiles() error {

	jobsQueue := make(chan pieceJob, len(files.torrent.PiecesHash))
	result := make(chan pieceResult)

	for _, peers := range files.peerClients {
		peers := peers

		go func() {
			for piece := range jobsQueue {

				pieceBuffer, err := peers.GetPiece(piece.index, piece.length, piece.hash)

				if err == nil {
					//fmt.Println("Downloaded piece number", piece.index)

				}
				if err != nil {
					jobsQueue <- piece
					if errors.Is(err, torrent.PieceNotFound) {
						continue
					}

					return
				}

				result <- pieceResult{
					Index:     piece.index,
					FilePiece: pieceBuffer,
				}

			}

		}()

	}

	for index, hash := range files.torrent.PiecesHash {

		length := files.torrent.PieceLength

		if index == len(files.torrent.PiecesHash)-1 {
			length = files.torrent.TotalLength - files.torrent.PieceLength*(len(files.torrent.PiecesHash)-1)

		}

		jobsQueue <- pieceJob{
			index:  index,
			length: length,
			hash:   hash,
		}

	}
	completedFilesBuffer := make([]byte, files.torrent.TotalLength)

	for count := 0; count < len(files.torrent.PiecesHash); count++ {

		piece := <-result

		copy(completedFilesBuffer[piece.Index*files.torrent.PieceLength:], piece.FilePiece)

		GBdone := (float64(count) / float64(len(files.torrent.PiecesHash))) * (float64(files.torrent.Files[0].Length) / (1024 * 1024 * 1024))
		fmt.Printf("%0.2f%% Completed: %.2f/%.2fGB \n", (float64(count) / float64(len(files.torrent.PiecesHash)) * 100), GBdone, float64(files.torrent.Files[0].Length)/(1024*1024*1024))

	}

	close(jobsQueue)

	fmt.Println("Writting Data to file")

	var currentFileByte int
	for _, file := range files.torrent.Files {
		outputPath := filepath.Join("D:/Anime", file.Path)

		filesBytes := completedFilesBuffer[currentFileByte : currentFileByte+file.Length]

		err := os.WriteFile(outputPath, filesBytes, os.ModePerm)
		currentFileByte += file.Length

		if err != nil {
			fmt.Println("File writting error", err)
		}

	}
	return nil
}
