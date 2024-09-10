package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunDecode(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "5:hello"})

	assert.Equal(t, buffer.String(), "\"hello\"\n")
	assert.Nil(t, err)

	buffer.Reset()
	err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "i52e"})
	assert.Equal(t, buffer.String(), "52\n")
	assert.Nil(t, err)

	buffer.Reset()
	err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "l5:helloi52ee"})
	assert.Equal(t, buffer.String(), "[\"hello\",52]\n")
	assert.Nil(t, err)

	buffer.Reset()
	err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "d3:foo3:bar5:helloi52ee"})
	assert.Equal(t, buffer.String(), "{\"foo\":\"bar\",\"hello\":52}\n")
	assert.Nil(t, err)

	buffer.Reset()
	err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "d4:spaml1:a1:bee"})
	assert.Equal(t, buffer.String(), "{\"spam\":[\"a\",\"b\"]}\n")
	assert.Nil(t, err)
}

func TestRunInfo(t *testing.T) {
	buffer := &bytes.Buffer{}
	NewBittorrentClient(&Config{Out: buffer}).Run([]string{"info", "../../sample.torrent"})

	expectedOutput := `Tracker URL: http://bittorrent-test-tracker.codecrafters.io/announce
Length: 92063
Info Hash: d69f91e6b2ae4c542468d1073a71d4ea13879a7f
Piece Length: 32768
Piece Hashes:
e876f67a2a8886e8f36b136726c30fa29703022d
6e2275e604a0766656736e81ff10b55204ad8d35
f00d937a0213df1982bc8d097227ad9e909acc17
`

	actual := buffer.String()
	assert.Equal(t, actual, expectedOutput)

}
func TestRunPeers(t *testing.T) {
	buffer := &bytes.Buffer{}
	NewBittorrentClient(&Config{Out: buffer}).Run([]string{"peers", "../../sample.torrent"})

	expectedOutput := `165.232.111.122:51437
161.35.47.237:51419
139.59.169.165:51487`
	actual := buffer.String()
	assert.Equal(t, actual, expectedOutput)

}

func TestRun_Invalid(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := NewBittorrentClient(&Config{Out: buffer}).Run([]string{})
	assert.NotNil(t, err)

	err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"unknown_command"})
	assert.NotNil(t, err)

	err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "invalid"})
	assert.NotNil(t, err)

	err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "-5:hmm"})
	assert.NotNil(t, err)

	err = NewBittorrentClient(&Config{Out: buffer}).Run([]string{"decode", "hi:"})
	assert.NotNil(t, err)

}
