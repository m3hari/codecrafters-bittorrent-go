package main

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
)

type TorrentInfo struct {
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
	PieceLength int    `bencode:"piece length"`
	// concatenated SHA-1 hashes of each piece (20 bytes each)
	Pieces string `bencode:"pieces"`
}

type Torrent struct {
	Announce string `bencode:"announce"`
	Info     *TorrentInfo
}

func New(fileName string) (torrent *Torrent, err error) {
	rawData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	result, err := bencode.Unmarshal(string(rawData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode bencode: %v", err)
	}

	data, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid torrent file. Root element is not a dictionary")
	}

	info := data["info"].(map[string]interface{})

	torrentInfo := &TorrentInfo{
		Name:        info["name"].(string),
		Length:      info["length"].(int),
		PieceLength: info["piece length"].(int),
		Pieces:      info["pieces"].(string),
	}

	torrent = &Torrent{
		Announce: data["announce"].(string),
		Info:     torrentInfo,
	}

	return torrent, nil
}

func (torrent *Torrent) InfoHash() ([20]byte, error) {
	bencodeData, err := bencode.ToBencodeDictionary(*torrent.Info)
	empty := [20]byte{}
	if err != nil {
		return empty, err
	}
	bencodedString, err := bencode.Marshal(bencodeData)
	if err != nil {
		return empty, fmt.Errorf("failed to bencode info: %v", err)
	}

	return sha1.Sum([]byte(bencodedString)), nil

}

func (torrent *Torrent) PieceHashes() ([]string, error) {
	pieces := torrent.Info.Pieces

	if len(pieces)%20 != 0 {
		return nil, fmt.Errorf("invalid pieces length")
	}

	result := make([]string, len(pieces)/20)
	for i := 0; i < len(pieces); i += 20 {
		result = append(result, fmt.Sprintf("%x\n", pieces[i:i+20]))
	}

	return result, nil
}
