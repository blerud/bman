package server

const (
	bombWidth       = 1
	bombHeight      = 1
	bombTimerLength = 120
)

type Bomb struct {
	entity *Entity
	state  uint8
}

func newBomb(id int32, x float32, y float32) *Entity {
	bomb := Bomb{
		nil,
		bombTimerLength,
	}
	entity := Entity{
		id,
		BOMB,
		x,
		y,
		bombWidth,
		bombHeight,
		&bomb,
	}
	bomb.entity = &entity
	return &entity
}

func (bomb *Bomb) encode() []byte {
	buffer := make([]byte, 14)
	buffer[0] = BOMB
	writeInt32ToBuffer(bomb.entity.entityId, buffer[1:])
	writeFloat32ToBuffer(bomb.entity.x, buffer[5:])
	writeFloat32ToBuffer(bomb.entity.y, buffer[9:])
	buffer[13] = bomb.state
	return buffer
}

func (bomb *Bomb) step(view EntitiesView) bool {
	return false
}
