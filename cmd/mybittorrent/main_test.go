package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunDecode(t *testing.T) {
	buffer := &bytes.Buffer{}
	result, err := NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "5:hello"})

	assert.Equal(t, result, "\"hello\"")
	assert.Nil(t, err)

	result, err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "i52e"})
	assert.Equal(t, result, "52")
	assert.Nil(t, err)

	result, err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "l5:helloi52ee"})
	assert.Equal(t, result, "[\"hello\",52]")
	assert.Nil(t, err)

	result, err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "d3:foo3:bar5:helloi52ee"})
	assert.Equal(t, result, "{\"foo\":\"bar\",\"hello\":52}")
	assert.Nil(t, err)

	result, err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "d4:spaml1:a1:bee"})
	assert.Equal(t, result, "{\"spam\":[\"a\",\"b\"]}")
	assert.Nil(t, err)
}

func TestRunInfo(t *testing.T) {
	buffer := &bytes.Buffer{}

	NewBittorrentClient(&Config{Out: buffer}).Run([]string{"info", "../../sample.torrent"})

	line, _ := buffer.ReadString('\n')
	assert.Equal(t, string(line), "Tracker URL: http://bittorrent-test-tracker.codecrafters.io/announce\n")

	line, _ = buffer.ReadString('\n')
	assert.Equal(t, string(line), "Length: 92063\n")

	line, _ = buffer.ReadString('\n')
	assert.Equal(t, string(line), "Info Hash: d69f91e6b2ae4c542468d1073a71d4ea13879a7f\n")

	line, _ = buffer.ReadString('\n')
	assert.Equal(t, string(line), "Piece Length: 32768\n")

	line, _ = buffer.ReadString('\n')
	assert.Equal(t, string(line), "Piece Hashes:\n")
	line, _ = buffer.ReadString('\n')
	assert.Equal(t, string(line), "e876f67a2a8886e8f36b136726c30fa29703022d\n")
	line, _ = buffer.ReadString('\n')
	assert.Equal(t, string(line), "6e2275e604a0766656736e81ff10b55204ad8d35\n")
	line, _ = buffer.ReadString('\n')
	assert.Equal(t, string(line), "f00d937a0213df1982bc8d097227ad9e909acc17\n")

}

func TestRun_Invalid(t *testing.T) {
	buffer := &bytes.Buffer{}
	_, err := NewBittorrentClient(&Config{Out: buffer}).Run([]string{})
	assert.NotNil(t, err)

	_, err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"unknown_command"})
	assert.NotNil(t, err)

	_, err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "invalid"})
	assert.NotNil(t, err)

	_, err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "-5:hmm"})
	assert.NotNil(t, err)

	_, err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "hi:"})
	assert.NotNil(t, err)

}
