package server

import (
	"encoding/binary"
	"math"
)

const (
	headerLength = 13

	messageServerHeartbeat = 0
	messageClientHeartbeat = 1
	messageGameInit = 2
	messageGameStart = 3
	messageCreated = 10
	messageUpdated = 11
	messageDeleted = 12
	messageClientAction = 20

	playerMessageLength = 15
)

type Message struct {
	messageType byte
	length int
	timestamp int64
	content []byte
}

type EntityMessage struct {
	entityType byte
	entityId int32
	x float32
	y float32
}

func writeInt32ToBuffer(i int32, buffer []byte) {
	binary.BigEndian.PutUint32(buffer, uint32(i))
}

func writeInt64ToBuffer(i int64, buffer []byte) {
	binary.BigEndian.PutUint64(buffer, uint64(i))
}

func writeFloat32ToBuffer(f float32, buffer []byte) {
	bits := math.Float32bits(f)
	binary.BigEndian.PutUint32(buffer, bits)
}

func readInt32FromBuffer(buffer []byte) int32 {
	return int32(binary.BigEndian.Uint32(buffer))
}

func readInt64FromBuffer(buffer []byte) int64 {
	return int64(binary.BigEndian.Uint64(buffer))
}

func readFloat32FromBuffer(buffer []byte) float32 {
	bits := binary.BigEndian.Uint32(buffer)
	return math.Float32frombits(bits)
}

func (m *Message) encodeMessage() []byte {
	message := make([]byte, len(m.content) + headerLength)
	message[0] = m.messageType
	writeInt32ToBuffer(int32(m.length), message[1:])
	writeInt64ToBuffer(m.timestamp, message[5:])
	for i, b := range m.content {
		message[i + headerLength] = b
	}
	return message
}
