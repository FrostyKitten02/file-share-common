package util

import (
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
	"net"
)

//TODO: Extract parsing logic, do the same in writter

func ReadConnectionInfo(conn net.Conn) (RoomConnectionInfo, error) {
	roomConnectionInfo := RoomConnectionInfo{}
	data, header, err := ReadPacket(conn)
	if err != nil {
		return roomConnectionInfo, err
	}

	if header.PacketType != PACKET_TYPE_ROOM_CONNECTION_INFO {
		return roomConnectionInfo, errors.New("wrong packet type")
	}

	roomConnectionInfo.Parse(data)
	return roomConnectionInfo, nil
}

func ReadConnectResponse(conn net.Conn) (ConnectResponse, error) {
	connectResponse := ConnectResponse{}
	data, header, err := ReadPacket(conn)
	if err != nil {
		return connectResponse, err
	}

	if header.PacketType != PACKET_TYPE_CONNECT_RESPONSE {
		return connectResponse, errors.New("wrong packet type")
	}

	connectResponse.Parse(data)
	return connectResponse, nil
}

func ReadConnectRequest(conn net.Conn) (ConnectRequest, error) {
	connectRequest := ConnectRequest{}
	data, header, err := ReadPacket(conn)
	if err != nil {
		return connectRequest, err
	}

	if header.PacketType != PACKET_TYPE_CONNECT_REQUEST {
		return connectRequest, errors.New("wrong packet type")
	}

	connectRequest.Parse(data)
	return connectRequest, nil
}

func ReadConnectToRoomMessage(conn net.Conn) (string, error) {
	data, header, err := ReadPacket(conn)
	if err != nil {
		return "", err
	}

	if header.PacketType != PACKET_TYPE_CONNECT_TO_ROOM {
		return "", errors.New("wrong packet type")
	}

	return string(data), nil
}

func ParseConnectToRoomMessage(data []byte) string {
	return string(data)
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
