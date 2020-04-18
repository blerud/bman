package server

const (
	PLAYER = iota
	BOMB
	ITEM
	WALL
)

const (
	BOMB_ITEM = iota
	FIRE
	MAX_FIRE
	KICK
	VEST
)

type EntityInfo interface {
	encode() []byte
}

type Entity struct {
	// x and y of bottom left corner of object -- all objects are square
	entityId   int32
	entityType byte
	x          float32
	y          float32
	width      float32
	height     float32

	entityInfo EntityInfo
}

func (e *Entity) encode() []byte {
	return e.entityInfo.encode()
}
