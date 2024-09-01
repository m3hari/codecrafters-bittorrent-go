package main

import (
	"fmt"
	"os"
)

type Info struct {
	name   string
	length int
}

type TorrentMetaInfo struct {
	announce string
	info     Info
}

func getTorrentMetaInfo(fileName string) (*TorrentMetaInfo, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	result, _, err := decodeBencode(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode bencode: %v", err)
	}

	info := TorrentMetaInfo{
		announce: result.(map[string]any)["announce"].(string),
		info: Info{
			name:   result.(map[string]any)["info"].(map[string]any)["name"].(string),
			length: result.(map[string]any)["info"].(map[string]any)["length"].(int),
		},
	}

	return &info, nil
}
