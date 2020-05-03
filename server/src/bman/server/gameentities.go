package server

const (
	PLAYER = iota
	BOMB
	ITEM
	HARD_WALL
	SOFT_WALL
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
	step(view EntitiesView) bool
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
