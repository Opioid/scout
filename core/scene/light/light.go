package light

import (
	"github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/base/math"
)

type Light interface {
	Entity() *entity.Entity
	
	Vector(p math.Vector3) math.Vector3
	Light(p, color math.Vector3) math.Vector3

	SetColor(color math.Vector3)
	SetLumen(lumen float32)
}

type light struct {
	entity entity.Entity
	color math.Vector3
	lumen float32
}

func (l *light) Entity() *entity.Entity {
	return &l.entity
}

func (l *light) SetColor(color math.Vector3) {
	l.color = color
}

func (l *light) SetLumen(lumen float32) {
	l.lumen = lumen
}