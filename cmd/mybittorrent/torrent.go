package main

import (
	"crypto/sha1"
	"fmt"
	"os"
	"reflect"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
)

type TorrentInfo struct {
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}
type Torrent struct {
	Announce string `bencode:"announce"`
	Info     *TorrentInfo
}

func NewTorrent(fileName string) (torrent *Torrent, err error) {
	rawData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	result, _, err := bencode.Unmarshal(string(rawData))
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

func ToBencodeDictionary(torrent any) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	value := reflect.ValueOf(torrent)
	typ := reflect.TypeOf(torrent)

	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Tag.Get("bencode")
		if fieldName == "" {
			fieldName = field.Name
		}
		fieldValue := value.Field(i).Interface()

		result[fieldName] = fieldValue
	}

	return result, nil
}

func InfoHash(info TorrentInfo) (string, error) {
	bencodeData, err := ToBencodeDictionary(info)
	if err != nil {
		return "", err
	}
	bencodedString, err := bencode.Marshal(bencodeData)
	if err != nil {
		return "", fmt.Errorf("failed to bencode info: %v", err)
	}

	return fmt.Sprintf("%x", sha1.Sum([]byte(bencodedString))), nil
}

func PiecesHashes(torrent Torrent) ([]string, error) {
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
