package texture

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	_ "fmt"
)

type addressMode interface {
	address2D(uv math.Vector2) math.Vector2
}

type AddressMode_clamp struct {
}

func (a *AddressMode_clamp) address2D(uv math.Vector2) math.Vector2 {
	u := math32.Clamp(uv.X, 0, 1)
	v := math32.Clamp(uv.Y, 0, 1)

	return math.MakeVector2(u, v)
}

type AddressMode_repeat struct {
}

func (a *AddressMode_repeat) address2D(uv math.Vector2) math.Vector2 {
	u := uv.X - math32.Floor(uv.X)
	v := uv.Y - math32.Floor(uv.Y)

	return math.MakeVector2(u, v)
}
