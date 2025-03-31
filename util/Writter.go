package util

import (
	"encoding/binary"
	"github.com/vmihailenco/msgpack/v5"
	"log/slog"
	"net"
)

// TODO add functionality where we can sand plain messages not only structs, similar to how len is sent!!
func WriteMessage(conn net.Conn, msg *Message) error {
	serialized, serErr := msgpack.Marshal(&msg)
	if serErr != nil {
		slog.Error("Error serializing message ", serErr.Error())
		return serErr
	}

	serializedLen := uint64(len(serialized))
	lenMessage := make([]byte, 8)
	binary.LittleEndian.PutUint64(lenMessage, serializedLen)

	data := append(lenMessage, serialized...)
	_, err := conn.Write(data)
	if err != nil {
		slog.Info("Write error:", err)
		return err
	}

	return nil
}
