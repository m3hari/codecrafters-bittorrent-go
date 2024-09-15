package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/codecrafters-io/bittorrent-starter-go/internal/bencode"
)

const peerId string = "00112233445566778899"
const peerPort string = "6881"

type Torrent struct {
	Announce string `bencode:"announce"`
	Info     *Info
}

type Info struct {
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
	PieceLength int    `bencode:"piece length"`
	// concatenated SHA-1 hashes of each piece (20 bytes each)
	Pieces string `bencode:"pieces"`
}

func NewTorrent(fileName string) (torrent *Torrent, err error) {
	rawData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
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

	torrentInfo := &Info{
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

func (t *Torrent) Handshake(peerAddress string) (peerId []byte, err error) {
	infoHash, err := t.InfoHash()
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	// Peer protocol handshake message is 68 bytes.
	var handshake = make([]byte, 68)
	handshake[0] = 19
	copy(handshake[1:20], []byte("BitTorrent protocol"))
	copy(handshake[28:48], infoHash[:])
	copy(handshake[48:], []byte(peerId))

	_, err = conn.Write(handshake)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 68)
	_, err = io.ReadFull(conn, response)
	if err != nil {
		return nil, err
	}

	return response[48:], nil
}

type TrackerResponse struct {
	Interval int
	Peers    []string
}

func (t *Torrent) DiscoverPeers() (*TrackerResponse, error) {
	infoHash, err := t.InfoHash()
	if err != nil {
		return nil, err
	}

	infoHashEscaped := url.QueryEscape(string(infoHash[:]))
	queryParams := url.Values{
		"peer_id":    {peerId},
		"port":       {peerPort},
		"uploaded":   {"0"},
		"downloaded": {"0"},
		"left":       {strconv.Itoa(t.Info.Length)},
		"compact":    {string("1")},
	}

	// encoding the info hash along with the other query params breaks the url
	resp, err := http.Get(t.Announce + "?" + queryParams.Encode() + "&info_hash=" + infoHashEscaped)

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

func parsePeersResponse(byteString string) ([]string, error) {
	data := []byte(byteString)

	if len(data) == 0 {
		return nil, fmt.Errorf("invalid response: empty peers field")
	}

	if len(data)%6 != 0 {
		return nil, fmt.Errorf("invalid response: peers field length (%d) is not a multiple of 6", len(data))
	}

	peers := make([]string, 0, len(data)/6)
	for i := 0; i < len(data); i += 6 {
		ip := net.IP(data[i : i+4])
		port := binary.BigEndian.Uint16(data[i+4 : i+6])
		peers = append(peers, fmt.Sprintf("%s:%d", ip.String(), port))
	}

	return peers, nil
}
