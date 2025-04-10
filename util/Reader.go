package util

import (
	"encoding/binary"
	"errors"
	"github.com/vmihailenco/msgpack/v5"
	"log/slog"
	"net"
)

type Message struct {
	Data string
}

func readBytes(conn net.Conn) ([]byte, error) {
	messageLength, lenErr := readLength(conn)

	if lenErr != nil {
		return nil, lenErr
	}

	maxReadLen := uint64(4096)
	readLen := maxReadLen
	if messageLength < maxReadLen {
		readLen = messageLength
	}

	buf := make([]byte, messageLength)
	tmp := make([]byte, readLen)

	readSize := uint64(0)
	for readSize < messageLength {
		//TODO: maybe use io.ReadFull instead, conn.Read may not return requested size but what ever is has in buff, io.ReadFull might return EOF err when done
		//this is so we don't have to handle checking if everything has been read
		//but we have to make sure we don't try to read more than it was requested
		n, err := conn.Read(tmp)
		if err != nil {
			return nil, err
		}
		buf = append(buf[:min(readSize-1, 0)], tmp[:n]...)
		readSize += uint64(n)
	}

	return buf, nil
}

func readLength(conn net.Conn) (uint64, error) {
	buf := make([]byte, 8)
	readLen, err := conn.Read(buf)
	if err != nil {
		return uint64(0), err
	}

	if readLen != 8 {
		return uint64(0), errors.New("Error reading message size")
	}

	return binary.LittleEndian.Uint64(buf), nil
}

// TODO: make more structured messages, using custom header and data.
// [header][data_len][data], header should contain packet type and any necessary information for processing
func ReadMessage(conn net.Conn) (*Message, error) {
	buf, err := readBytes(conn)
	if err != nil {
		slog.Info("Read error:", err)
		return nil, err
	}

	message := Message{}
	deserializeErr := msgpack.Unmarshal(buf, &message)
	if deserializeErr != nil {
		slog.Info("Deserialize error:", deserializeErr)
		return nil, deserializeErr
	}

	return &message, nil
}
