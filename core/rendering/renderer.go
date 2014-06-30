package rendering

import (
	"github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/camera"
	"github.com/Opioid/scout/base/math"
)

type Renderer struct {

}

func (r *Renderer) Render(scene *scene.Scene, context *Context) {
	target := context.Target
	dimensions := target.Dimensions()

	var ray math.Ray

	for y := 0; y < dimensions.Y; y++ {
		for x := 0; x < dimensions.X; x++ {

			context.Camera.GenerateRay(&camera.Sample{math.Vector2i{x, y}}, &ray)

			r := float32(y) / float32(dimensions.Y)
			g := float32(x) / float32(dimensions.X)

			var thit float32
			if scene.Intersect(&ray, &thit) {
				target.Set(x, y, math.Vector4{r, g, 0.5, 1.0})
			} else {
				target.Set(x, y, math.Vector4{0.0, 0.0, 0.0, 1.0})
			}
		}
	}
}