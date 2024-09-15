package main

import (
	"io"
	"net"
)

func handshake(torrent *Torrent, peerAddress string) (string, error) {
	infoHash, err := torrent.InfoHash()
	if err != nil {
		return "", err
	}

	connection, err := net.Dial("tcp", peerAddress)
	if err != nil {
		return "", err
	}

	defer func() {
		err = connection.Close()
	}()

	// Based on peer protocol a handshake message 68 bytes long with different
	// parts. See: https://www.bittorrent.org/beps/bep_0003.html#peer-protocol
	var message = []byte{}
	message = append(message, 19)
	message = append(message, []byte("BitTorrent protocol")...)
	message = append(message, make([]byte, 8)...)
	message = append(message, infoHash[:]...)
	message = append(message, []byte("00112233445566778899")...)

	_, err = connection.Write(message)
	if err != nil {
		return "", err
	}

	reply := make([]byte, 68)
	_, err = io.ReadFull(connection, reply)
	if err != nil {
		return "", err
	}

	return string(reply), nil
}
