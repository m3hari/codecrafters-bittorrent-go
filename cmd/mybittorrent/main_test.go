package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunDecode(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"string", "5:hello", "\"hello\"\n"},
		{"integer", "i52e", "52\n"},
		{"list", "l5:helloi52ee", "[\"hello\",52]\n"},
		{"dictionary", "d3:foo3:bar5:helloi52ee", "{\"foo\":\"bar\",\"hello\":52}\n"},
		{"nested", "d4:spaml1:a1:bee", "{\"spam\":[\"a\",\"b\"]}\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			client := NewClient(buffer)
			err := client.Run([]string{"decode", tc.input})

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, buffer.String())
		})
	}
}

func TestRunInfo(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := NewClient(buffer).Run([]string{"info", "../../sample.torrent"})

	require.NoError(t, err)

	expectedOutput := `Tracker URL: http://bittorrent-test-tracker.codecrafters.io/announce
Length: 92063
Info Hash: d69f91e6b2ae4c542468d1073a71d4ea13879a7f
Piece Length: 32768
Piece Hashes:
e876f67a2a8886e8f36b136726c30fa29703022d
6e2275e604a0766656736e81ff10b55204ad8d35
f00d937a0213df1982bc8d097227ad9e909acc17
`

	assert.Equal(t, expectedOutput, buffer.String())
}

func TestRunPeers(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := NewClient(buffer).Run([]string{"peers", "../../sample.torrent"})

	require.NoError(t, err)

	expectedOutput := `165.232.111.122:51437
161.35.47.237:51419
139.59.169.165:51487`

	assert.Equal(t, expectedOutput, buffer.String())
}

func TestRunInvalidCases(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{"no arguments", []string{}},
		{"unknown command", []string{"unknown_command"}},
		{"invalid decode input", []string{"decode", "invalid"}},
		{"negative string length", []string{"decode", "-5:hmm"}},
		{"invalid string format", []string{"decode", "hi:"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			client := NewClient(buffer)
			err := client.Run(tc.args)

			assert.Error(t, err)
		})
	}
}

func TestRunHandshake(t *testing.T) {
	// This test requires mocking the network connection
	// For simplicity, we'll just test the error case here
	buffer := &bytes.Buffer{}
	err := NewClient(buffer).Run([]string{"handshake", "nonexistent.torrent", "127.0.0.1:6881"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}

func TestRunUnknownCommand(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := NewClient(buffer).Run([]string{"unknown"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown command")
}
