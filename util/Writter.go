package util

import (
	"encoding/binary"
	"net"
)

func WriteStringMessage(conn net.Conn, msg string) error {
	data := []byte(msg)
	packet := createPacket(PACKET_TYPE_TEXT, data)

	//TODO check if all was written!
	_, err := conn.Write(packet)
	if err != nil {
		return err
	}

	return nil
}

func createPacketHeader(packetType uint8, dataLen uint64) []byte {
	header := make([]byte, 9)
	header[0] = packetType
	binary.LittleEndian.PutUint64(header[1:], dataLen)
	return header
}

func createPacket(packetType uint8, data []byte) []byte {
	header := createPacketHeader(packetType, uint64(len(data)))
	packet := make([]byte, len(data)+len(header))
	copy(packet[0:9], header)
	copy(packet[9:], data)
	return packet
}
