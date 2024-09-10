package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
)

type TrackerResponse struct {
	Interval int
	Peers    []string
}

func DiscoverPeers(torrent *Torrent) (*TrackerResponse, error) {
	infoHash, err := torrent.InfoHash()
	if err != nil {
		return nil, err
	}

	infoHashEscaped := url.QueryEscape(string(infoHash[:]))
	queryParams := url.Values{
		"peer_id":    {"00112233445566778899"},
		"port":       {"6881"},
		"uploaded":   {"0"},
		"downloaded": {"0"},
		"left":       {strconv.Itoa(torrent.Info.Length)},
		"compact":    {string("1")},
	}

	// encoding the info hash along with the other query params breaks the url
	resp, err := http.Get(torrent.Announce + "?" + queryParams.Encode() + "&info_hash=" + infoHashEscaped)

	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch peers. HTTP status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body. err:%v", err)
	}

	rawBody := string(body)
	value, err := bencode.Unmarshal(rawBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response. response: %v err:%v", rawBody, err)
	}

	data, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response. root value is not a map. response: %v", rawBody)
	}

	interval, ok := data["interval"].(int)
	if !ok {
		return nil, fmt.Errorf("invalid response. Could not find interval. response: %v", rawBody)
	}

	peers, ok := data["peers"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid response.  Could not find peers. response: %v", rawBody)
	}

	return &TrackerResponse{
		Interval: interval,
		Peers:    parsePeersResponse(peers),
	}, nil
}

//	parsePeersResponse parses peers address from raw response
//
// In the compact representation of peers response peers field is of type
// string, but its content is binary data. Each peer is a sequence of 6 bytes.
func parsePeersResponse(byteString string) []string {
	data := []byte(byteString)

	peers := [][]byte{}
	start := 0
	for start+6 <= len(data) {
		end := start + 6
		peer := data[start:end]
		peers = append(peers, peer)
		start = end
	}

	result := []string{}
	// for each peer , parse address
	for _, peer := range peers {
		address := fmt.Sprintf("%d.%d.%d.%d:%d",
			peer[0],
			peer[1],
			peer[2],
			peer[3],
			binary.BigEndian.Uint16(peer[4:6]),
		)
		result = append(result, address)
	}

	return result
}
