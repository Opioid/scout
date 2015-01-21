package light

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Light interface {
	Prop() *prop.Prop

	SetColor(color math.Vector3)
	SetLumen(lumen float32)

	Samples(p math.Vector3, time float32, subsample, maxSamples uint32, sampler sampler.Sampler, samples []Sample) []Sample
}

type light struct {
	prop prop.Prop
	color math.Vector3
	lumen float32
}

func (l *light) Prop() *prop.Prop {
	return &l.prop
}

func (l *light) SetColor(color math.Vector3) {
	l.color = color
}

func (l *light) SetLumen(lumen float32) {
	l.lumen = lumen
}