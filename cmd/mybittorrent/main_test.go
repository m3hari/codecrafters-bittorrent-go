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
