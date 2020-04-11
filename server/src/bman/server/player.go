package server

const (
	createMessageLength = 15
	updateMessageLength = 14
	deleteMessageLength = 14
)

const (
	playerWidth = float32(0.5)
	playerHeight = float32(0.5)
	directionUp = 0
	directionLeft = 1
	directionDown = 2
	directionRight = 3
)

type Player struct {
	entity *Entity
	action byte
	direction byte
}

type PlayerMessage struct {
	EntityMessage
	action byte
	direction byte
}

func newPlayer(id int32, x float32, y float32) *Entity {
	player := Player {
		nil,
		0,
		directionUp,
	}
	entity := Entity{
		id,
		PLAYER,
		x,
		y,
		playerWidth,
		playerHeight,
		&player,
	}
	player.entity = &entity
	return &entity
}

func (p *Player) processUpdate(bytes []byte) int {
	playerMessage := decode(bytes)
	p.update(playerMessage)
	return updateMessageLength
}

func (p *Player) encode() []byte {
	playerUpdate := PlayerMessage{
		EntityMessage{
			p.entity.entityType,
			p.entity.entityId,
			p.entity.x,
			p.entity.y,
		},
		p.action,
		p.direction,
	}
	buffer := make([]byte, 14)
	buffer[0] = playerUpdate.entityType
	writeInt32ToBuffer(playerUpdate.entityId, buffer[1:])
	writeFloat32ToBuffer(playerUpdate.x, buffer[5:])
	writeFloat32ToBuffer(playerUpdate.y, buffer[9:])
	buffer[13] = playerUpdate.action
	buffer[14] = playerUpdate.direction
	return buffer
}

func (p *Player) update(message PlayerMessage) {
	p.entity.x = message.x
	p.entity.y = message.y
	p.action = message.action
	p.direction = message.direction
}

func decode(content []byte) PlayerMessage {
	entityType := content[0]
	entityId := readInt32FromBuffer(content[1:])
	x := readFloat32FromBuffer(content[5:])
	y := readFloat32FromBuffer(content[8:])
	action := content[13]
	direction := content[14]
	return PlayerMessage{
		EntityMessage{
			entityType,
			entityId,
			x,
			y,
		},
		action,
		direction,
	}
}
