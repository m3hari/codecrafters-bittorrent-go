package main

import (
	"encoding/json"
	"errors"
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

func (client *BittorrentClient) Run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: <command> <argument>")
	}

	command := args[0]
	switch {
	case command == "decode":
		result, err := bencode.Unmarshal(args[1])
		if err != nil {
			return err
		}
		jsonOutput, err := json.Marshal(result)
		if err != nil {
			return err
		}
		client.Out.Write((jsonOutput))
		client.Out.Write(([]byte("\n")))
		return nil

	case command == "info":
		torrent, err := New(args[1])
		if err != nil {
			return err
		}

		infoHash, err := torrent.InfoHash()
		if err != nil {
			return err
		}

		pieceHashes, err := torrent.PieceHashes()
		if err != nil {
			return err
		}

		client.Out.Write([]byte(fmt.Sprintf("Tracker URL: %v\n", torrent.Announce)))
		client.Out.Write([]byte(fmt.Sprintf("Length: %v\n", torrent.Info.Length)))
		client.Out.Write([]byte(fmt.Sprintf("Info Hash: %v\n", fmt.Sprintf("%x", infoHash))))
		client.Out.Write([]byte(fmt.Sprintf("Piece Length: %v\n", torrent.Info.PieceLength)))
		client.Out.Write([]byte("Piece Hashes:\n"))
		for _, item := range pieceHashes {
			client.Out.Write([]byte(item))
		}

		return nil
	default:
		return errors.ErrUnsupported
	}
}

func main() {
	client := NewBittorrentClient(&Config{})

	err := client.Run(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
}
