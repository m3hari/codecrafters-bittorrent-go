package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
)

type Client struct {
	out io.Writer
}

func NewClient(out io.Writer) *Client {
	if out == nil {
		out = os.Stdout
	}
	return &Client{out: out}
}

func (c *Client) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: <command> <argument>")
	}

	cmd := args[0]
	handler, ok := commandHandlers[cmd]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd)
	}

	return handler(c, args[1:])
}

var commandHandlers = map[string]func(*Client, []string) error{
	"decode":    decodeCommand,
	"info":      infoCommand,
	"peers":     peersCommand,
	"handshake": handshakeCommand,
}

func decodeCommand(c *Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: decode <bencoded string>")
	}
	result, err := bencode.Unmarshal(args[0])
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}
	return json.NewEncoder(c.out).Encode(result)
}

func infoCommand(c *Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: info <torrent file>")
	}
	torrent, err := NewTorrent(args[0])
	if err != nil {
		return fmt.Errorf("failed to create torrent: %w", err)
	}

	infoHash, err := torrent.InfoHash()
	if err != nil {
		return fmt.Errorf("failed to get info hash: %w", err)
	}

	pieceHashes, err := torrent.PieceHashes()
	if err != nil {
		return fmt.Errorf("failed to get piece hashes: %w", err)
	}

	fmt.Fprintf(c.out, "Tracker URL: %s\n", torrent.Announce)
	fmt.Fprintf(c.out, "Length: %d\n", torrent.Info.Length)
	fmt.Fprintf(c.out, "Info Hash: %x\n", infoHash)
	fmt.Fprintf(c.out, "Piece Length: %d\n", torrent.Info.PieceLength)
	fmt.Fprintln(c.out, "Piece Hashes:")
	for _, hash := range pieceHashes {
		fmt.Fprintf(c.out, "%s", hash)
	}

	return nil
}

func peersCommand(c *Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: peers <torrent file>")
	}
	torrent, err := NewTorrent(args[0])
	if err != nil {
		return fmt.Errorf("failed to create torrent: %w", err)
	}
	result, err := torrent.DiscoverPeers()
	if err != nil {
		return fmt.Errorf("failed to discover peers: %w", err)
	}

	fmt.Fprintf(c.out, "%s", strings.Join(result.Peers, "\n"))

	return nil
}

func handshakeCommand(c *Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: handshake <torrent file> <peer address>")
	}

	torrent, err := NewTorrent(args[0])
	if err != nil {
		return fmt.Errorf("failed to create torrent: %w", err)
	}

	peerAddress := args[1]
	if _, _, err := net.SplitHostPort(peerAddress); err != nil {
		return fmt.Errorf("invalid peer address: %w", err)
	}

	peerID, err := torrent.Handshake(args[1])
	if err != nil {
		return fmt.Errorf("handshake failed: %w", err)
	}

	fmt.Fprintf(c.out, "Peer ID: %x\n", peerID)
	return nil
}

func main() {
	client := NewClient(nil)
	if err := client.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
