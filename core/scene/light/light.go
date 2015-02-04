package light

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Light interface {
	Prop() *prop.Prop

	Color() math.Vector3
	SetColor(color math.Vector3)

	Lumen() float32
	SetLumen(lumen float32)

	Sample(transformation *math.ComposedTransformation, p math.Vector3, time float32, subsample uint32, sampler sampler.Sampler) Sample
}

type light struct {
	prop prop.Prop
	color math.Vector3
	lumen float32
}

func (l *light) Prop() *prop.Prop {
	return &l.prop
}

func (l *light) Color() math.Vector3 {
	return l.color
}

func (l *light) SetColor(color math.Vector3) {
	l.color = color
}

func (l *light) Lumen() float32 {
	return l.lumen
}

func (l *light) SetLumen(lumen float32) {
	l.lumen = lumen
}

