package util

import (
	"encoding/binary"
	"net"
)

type PacketHeader struct {
	Len        uint64
	PacketType uint8
}

type ParsablePacket interface {
	Parse(data []byte)
	Serialize() []byte
}

type ConnectRequest struct {
	From string
}

func (c *ConnectRequest) Serialize() []byte {
	return []byte(c.From)
}

func (c *ConnectRequest) Parse(data []byte) {
	c.From = string(data)
}

type ConnectResponse struct {
	RoomConnectionInfo RoomConnectionInfo
	Allow              bool
	Id                 string
}

func (c *ConnectResponse) Serialize() []byte {
	roomConnection := c.RoomConnectionInfo.Serialize()
	roomConnectionLen := len(roomConnection)
	idLen := len(c.Id)
	data := make([]byte, 8+roomConnectionLen+1+idLen)

	binary.LittleEndian.PutUint64(data[0:], uint64(roomConnectionLen))

	offset := 9
	copy(data[offset:], roomConnection)

	offset = offset + roomConnectionLen
	allowByte := byte(0)
	if c.Allow {
		allowByte = byte(1)
	}
	copy(data[offset:], []byte{allowByte})

	offset = offset + 1
	copy(data[offset:], c.Id)

	return data
}

func (c *ConnectResponse) Parse(data []byte) {
	roomInfoLenBuf := data[0:9]
	roomInfoLen := int(binary.LittleEndian.Uint64(roomInfoLenBuf))

	offset := 9
	roomInfo := RoomConnectionInfo{}
	roomInfo.Parse(data[offset : offset+roomInfoLen])
	c.RoomConnectionInfo = roomInfo

	otherData := data[offset+roomInfoLen+1:]
	c.Allow = otherData[0] == 1
	c.Id = string(otherData[1:])
}

type RoomConnectionInfo struct {
	Ip   net.IP
	Port int
}

func (r *RoomConnectionInfo) Parse(data []byte) {
	portBuf := data[:8]
	ipBuf := data[8:]
	port := int(binary.BigEndian.Uint64(portBuf))
	r.Ip = net.ParseIP(string(ipBuf))
	r.Port = port
}

func (r *RoomConnectionInfo) Serialize() []byte {
	ipStr := r.Ip.String()
	ipBuf := []byte(ipStr)

	portBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(portBuf, uint64(r.Port))

	ipLen := len(ipBuf)
	portLen := len(portBuf)

	data := make([]byte, ipLen+portLen)
	copy(data[:portLen], portBuf)
	copy(data[portLen:], ipBuf)
	return data
}
