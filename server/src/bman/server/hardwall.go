package server

const (
	wallWidth  = 1
	wallHeight = 1
)

type HardWall struct {
	entity *Entity
}

func newHardWall(id int32, x float32, y float32) *Entity {
	hardWall := HardWall{
		nil,
	}
	entity := Entity{
		id,
		HARD_WALL,
		x,
		y,
		wallWidth,
		wallHeight,
		&hardWall,
	}
	hardWall.entity = &entity
	return &entity
}

func (wall *HardWall) encode() []byte {
	buffer := make([]byte, 13)
	buffer[0] = HARD_WALL
	writeInt32ToBuffer(wall.entity.entityId, buffer[1:])
	writeFloat32ToBuffer(wall.entity.x, buffer[5:])
	writeFloat32ToBuffer(wall.entity.y, buffer[9:])
	return buffer
}

func (wall *HardWall) step(view EntitiesView) bool {
	return false
}
