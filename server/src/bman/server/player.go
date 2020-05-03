package server

import (
	"math"
)

const (
	actionMessageLength = 5
)

const (
	playerWidth    = float32(1)
	playerHeight   = float32(1)
	playerXSpeed   = float32(3)
	playerYSpeed   = float32(3)
	directionUp    = 0
	directionLeft  = 1
	directionDown  = 2
	directionRight = 3
)

type Player struct {
	entity         *Entity
	action         PlayerAction
	direction      byte
	xSpeed         float32
	ySpeed         float32
	numBombs       int
	bombsAvailable int
	numFire        int
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
		0,
		0,
		1,
		1,
		1,
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

func (p *Player) step(view EntitiesView) bool {
	xSpeed := playerXSpeed / tick
	ySpeed := playerYSpeed / tick
	locX, locY := p.currentSquare()

	action := p.action

	if action.up {
		p.ySpeed = -ySpeed
	}
	if action.down {
		p.ySpeed = ySpeed
	}
	if action.left {
		p.xSpeed = -xSpeed
	}
	if action.right {
		p.xSpeed = xSpeed
	}
	newX := p.entity.x + p.xSpeed
	newY := p.entity.y + p.ySpeed
	didMove := p.xSpeed != 0 || p.ySpeed != 0

	collisions := view.collisions(p.entity, newX, newY)

	if action.bomb && p.bombsAvailable > 0 {
		// todo check if currently colliding with a bomb, if so don't place the bomb
		success := view.create(newBomb(view.genEntityId(), locX, locY, p.numFire))
		if success {
			p.bombsAvailable--
		}
	}

	if onSquare(newX, newY) {
		p.xSpeed = 0
		p.ySpeed = 0
	}

	if len(collisions) != 0 {
		p.xSpeed = 0
		p.ySpeed = 0
	}

	if didMove && len(collisions) == 0 {
		p.entity.x = newX
		p.entity.y = newY
		return true
	}

	return false
}

func (p *Player) update(action PlayerAction) {
	p.action = action
}

func (p *Player) currentSquare() (float32, float32) {
	return float32(math.Round(float64(p.entity.x))), float32(math.Round(float64(p.entity.y)))
}

func onSquare(x float32, y float32) bool {
	return wholeNumber(x) && wholeNumber(y)
}

func wholeNumber(val float32) bool {
	_, frac := math.Modf(float64(val))
	return frac < 1e-3 || frac > 1-1e-3
}

func decode(content []byte) PlayerAction {
	actionByte := content[4]
	return PlayerAction{
		up:    actionByte&1 == 1,
		down:  (actionByte>>1)&1 == 1,
		left:  (actionByte>>2)&1 == 1,
		right: (actionByte>>3)&1 == 1,
		bomb:  (actionByte>>4)&1 == 1,
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
