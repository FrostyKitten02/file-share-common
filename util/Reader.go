package util

import (
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"net"
)

type Message struct {
	Data string
}

// TODO: return packet struct containing header and data!
func ReadPacket(conn net.Conn) ([]byte, PacketHeader, error) {
	packetHeader, headerErr := readPacketHeader(conn)
	if headerErr != nil {
		return nil, packetHeader, headerErr
	}

	packetData := make([]byte, packetHeader.Len)

	_, err := io.ReadFull(conn, packetData)
	if err != nil {
		return nil, packetHeader, err
	}

	return packetData, packetHeader, nil
}

// packet header contains [type][data_len], type = 1byte, data_len = 8bytes
func readPacketHeader(conn net.Conn) (PacketHeader, error) {
	header := PacketHeader{}
	typeBuf := make([]byte, 1)
	_, typeErr := io.ReadFull(conn, typeBuf)
	if typeErr != nil {
		slog.Error(typeErr.Error())
		return PacketHeader{}, errors.New("error reading packet type")
	}
	header.PacketType = typeBuf[0]

	lenBuf := make([]byte, 8)
	_, lenErr := io.ReadFull(conn, lenBuf)
	if lenErr != nil {
		slog.Error(lenErr.Error())
		return header, errors.New("error reading packet length")
	}

	header.Len = binary.LittleEndian.Uint64(lenBuf)
	return header, nil
}
