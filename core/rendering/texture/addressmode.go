package texture

import (
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type addressMode interface {
	address2D(uv math.Vector2) math.Vector2
}

type AddressMode_clamp struct {
}

func (a *AddressMode_clamp) address2D(uv math.Vector2) math.Vector2 {
	u := math.Clamp(uv.X, 0.0, 1.0)
	v := math.Clamp(uv.Y, 0.0, 1.0)
	return math.Vector2{u, v}
}

type AddressMode_repeat struct {

}

func (a *AddressMode_repeat) address2D(uv math.Vector2) math.Vector2 {
	u := uv.X - math.Floor(uv.X)
	v := uv.Y - math.Floor(uv.Y)

//	fmt.Printf("%f -> %f \n", uv.X, u)

	return math.Vector2{u, v}
}