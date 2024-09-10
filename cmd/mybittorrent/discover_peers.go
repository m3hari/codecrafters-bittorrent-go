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

	peersList, err := parsePeersResponse(peers)
	if err != nil {
		return nil, err
	}

	return &TrackerResponse{
		Interval: interval,
		Peers:    peersList,
	}, nil
}

// parsePeersResponse parses the compact peer response into a list of peer addresses.
//
// Each peer is represented by 6 bytes: the first 4 bytes are the IP, and the last 2 bytes are the port.
func parsePeersResponse(byteString string) ([]string, error) {
	data := []byte(byteString)

	if len(data)%6 != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid response. peers field length is not a multiple of 6. response: %v", byteString)
	}

	peers := make([]string, 0, len(data)/6)
	for i := 0; i < len(data); i += 6 {
		ip := fmt.Sprintf("%d.%d.%d.%d", data[i], data[i+1], data[i+2], data[i+3])
		port := binary.BigEndian.Uint16(data[i+4 : i+6])
		peers = append(peers, fmt.Sprintf("%s:%d", ip, port))
	}

	return peers, nil
}
