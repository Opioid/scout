package rendering

import "github.com/Opioid/scout/base/math"

type Renderer struct {

}

func (r *Renderer) Render(target *PixelBuffer) {
	dimensions := target.Dimensions()
	for y := 0; y < dimensions.Y; y++ {
		for x := 0; x < dimensions.X; x++ {
			r := float32(y) / float32(dimensions.Y)
			g := float32(x) / float32(dimensions.X)

			target.Set(x, y, math.Vector4{r, g, 0.5, 1.0})
		}
	}
}