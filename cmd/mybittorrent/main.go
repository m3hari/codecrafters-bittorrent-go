package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

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
	if len(args) < 1 {
		return fmt.Errorf("usage: <command> <argument>")
	}

	command := args[0]
	switch {
	case command == "decode":
		if len(args) < 2 {
			return fmt.Errorf("usage: decode <bencoded string>")
		}
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
		if len(args) < 2 {
			return fmt.Errorf("usage: info <torrent file>")
		}
		torrent, err := NewTorrent(args[1])
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

	case command == "peers":
		if len(args) < 2 {
			return fmt.Errorf("usage: peers <torrent file>")
		}
		torrent, err := NewTorrent(args[1])
		if err != nil {
			return err
		}
		result, err := DiscoverPeers(torrent)
		if err != nil {
			return err
		}

		client.Out.Write([]byte(strings.Join(result.Peers, "\n")))

		return nil

	case command == "handshake":
		if len(args) < 3 {
			return fmt.Errorf("usage: handshake <torrent file> <peer address>")
		}

		torrent, err := NewTorrent(args[1])
		if err != nil {
			return err
		}

		reply, err := handshake(torrent, args[2])
		if err != nil {
			return err
		}

		client.Out.Write([]byte(fmt.Sprintf("Peer ID: %x\n", reply[48:])))

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
