package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
)

type Config struct {
	Args []string
	Out  io.Writer
}

type BittorrentClient struct {
	Out io.Writer
}

func NewBittorrentClient(cfg *Config) *BittorrentClient {
	client := &BittorrentClient{Out: os.Stdout}

	if cfg != nil && cfg.Out != nil {
		client.Out = cfg.Out
	}

	return client
}

func (client *BittorrentClient) Run(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("usage: <command> <argument>")
	}

	command := args[0]
	switch {
	case command == "decode":
		result, _, err := bencode.Unmarshal(args[1])
		if err != nil {
			return "", err
		}
		jsonOutput, err := json.Marshal(result)
		if err != nil {
			return "", err
		}
		return string(jsonOutput), nil

	case command == "info":
		torrent, err := NewTorrent(args[1])
		if err != nil {
			return "", err
		}

		infoHash, err := torrent.InfoHash()
		if err != nil {
			return "", err
		}

		pieceHashes, err := torrent.PieceHashes()
		if err != nil {
			return "", err
		}

		client.Out.Write([]byte(fmt.Sprintf("Tracker URL: %v\n", torrent.Announce)))
		client.Out.Write([]byte(fmt.Sprintf("Length: %v\n", torrent.Info.Length)))
		client.Out.Write([]byte(fmt.Sprintf("Info Hash: %v\n", infoHash)))
		client.Out.Write([]byte(fmt.Sprintf("Piece Length: %v\n", torrent.Info.PieceLength)))
		client.Out.Write([]byte("Piece Hashes:\n"))
		for _, item := range pieceHashes {
			client.Out.Write([]byte(item))
		}

		return "", nil

	default:
		return "", fmt.Errorf("Unknown command: " + command)
	}
}

func main() {
	client := NewBittorrentClient(&Config{})

	result, err := client.Run(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}
