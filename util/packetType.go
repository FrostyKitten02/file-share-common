package util

const (
	PACKET_TYPE_UNKNOWN      = 0
	PACKET_TYPE_TEXT         = 1 //simple text packet, more for information purposes
	PACKET_TYPE_ROOM_CREATED = 2 //when server creates a room for client

	PACKET_TYPE_CONNECT_TO_ROOM      = 3 //for when client wants to connect to other client room
	PACKET_TYPE_CONNECT_REQUEST      = 4 //used on server to ask client if they want to accept someone
	PACKET_TYPE_CONNECT_RESPONSE     = 5 //used on server, response from client from connection request, when someone wanted to connect to them
	PACKET_TYPE_ROOM_CONNECTION_INFO = 4 //connection info sent to client when he wants to connect to someone
)
