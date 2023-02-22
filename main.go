package main

import (
	"log"
	"os"

	"torrproject/torr"
)

func main() {
	inputPath := os.Args[1]
	outputPath := os.Args[2]

	torrent := torr.NewTorrent(inputPath, outputPath)

	log.Print("Analyzing input")
	_, err := torrent.ParseTorrent()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Locating peers peers")
	torrent.DiscoverPeers()

	log.Print("Initiating download")
	torrent.Download()

	log.Print("Writinging the file")
	torrent.OutputToFile()
}
