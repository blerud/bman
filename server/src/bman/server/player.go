package server

const (
	actionMessageLength = 5
)

const (
	playerWidth    = float32(0.5)
	playerHeight   = float32(0.5)
	playerXSpeed   = float32(1)
	playerYSpeed   = float32(1)
	directionUp    = 0
	directionLeft  = 1
	directionDown  = 2
	directionRight = 3
)

type Player struct {
	entity    *Entity
	action    PlayerAction
	direction byte
}

type PlayerMessage struct {
	EntityMessage
	action    byte
	direction byte
}

type PlayerAction struct {
	up    bool
	down  bool
	left  bool
	right bool
	bomb  bool
}

func newPlayer(id int32, x float32, y float32) *Entity {
	player := Player{
		nil,
		PlayerAction{
			false,
			false,
			false,
			false,
			false,
		},
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

func (p *Player) processPlayerAction(message Message) {
	playerAction := decode(message.content)
	p.update(playerAction)
}

func (p *Player) encode() []byte {
	playerUpdate := PlayerMessage{
		EntityMessage{
			p.entity.entityType,
			p.entity.entityId,
			p.entity.x,
			p.entity.y,
		},
		p.action.toByte(),
		p.direction,
	}
	buffer := make([]byte, 15)
	buffer[0] = playerUpdate.entityType
	writeInt32ToBuffer(playerUpdate.entityId, buffer[1:])
	writeFloat32ToBuffer(playerUpdate.x, buffer[5:])
	writeFloat32ToBuffer(playerUpdate.y, buffer[9:])
	buffer[13] = playerUpdate.action
	buffer[14] = playerUpdate.direction
	return buffer
}

func (p *Player) update(action PlayerAction) {
}

func decode(content []byte) PlayerAction {
	actionByte := content[4]
	return PlayerAction{
		up:    actionByte&1 == 1,
		down:  actionByte&(1<<1) == 1,
		left:  actionByte&(1<<2) == 1,
		right: actionByte&(1<<3) == 1,
		bomb:  actionByte&(1<<4) == 1,
	}
}

func (a *PlayerAction) toByte() byte {
	action := byte(0)
	if a.up {
		action |= 1
	}
	if a.down {
		action |= 1 << 1
	}
	if a.left {
		action |= 1 << 2
	}
	if a.right {
		action |= 1 << 3
	}
	if a.bomb {
		action |= 1 << 4
	}
	return action
}
