package util

import (
	"encoding/binary"
	"net"
)

func WriteRoomConnectionInfo(conn net.Conn, connInfo RoomConnectionInfo) error {
	data := connInfo.Serialize()
	packet := createPacket(PACKET_TYPE_ROOM_CONNECTION_INFO, data)

	_, err := conn.Write(packet)
	if err != nil {
		return err
	}

	return nil
}

func WriteConnectResponse(conn net.Conn, res ConnectResponse) error {
	data := res.Serialize()
	packet := createPacket(PACKET_TYPE_CONNECT_RESPONSE, data)

	_, err := conn.Write(packet)
	if err != nil {
		return err
	}

	return nil
}

func WriteConnectRequest(conn net.Conn, request ConnectRequest) error {
	packet := createPacket(PACKET_TYPE_CONNECT_REQUEST, request.Serialize())

	_, err := conn.Write(packet)
	if err != nil {
		return err
	}

	return nil
}

func WriteConnectToRoomMessage(conn net.Conn, roomId string) error {
	data := []byte(roomId)
	packet := createPacket(PACKET_TYPE_CONNECT_TO_ROOM, data)

	_, err := conn.Write(packet)
	if err != nil {
		return err
	}

	return nil
}

func WriteRoomCreatedMessage(conn net.Conn, roomId string) error {
	data := []byte(roomId)
	packet := createPacket(PACKET_TYPE_ROOM_CREATED, data)

	_, err := conn.Write(packet)
	if err != nil {
		return err
	}

	return nil
}

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
	binary.LittleEndian.PutUint64(header[1:], 18446744073709551615)
	return header
}

func createPacket(packetType uint8, data []byte) []byte {
	header := createPacketHeader(packetType, uint64(len(data)))
	packet := make([]byte, len(data)+len(header))
	copy(packet[0:9], header)
	copy(packet[9:], data)
	return packet
}
