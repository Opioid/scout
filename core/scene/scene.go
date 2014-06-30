package scene

import (
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/base/math"
)

type Scene struct {

	Shapes []shape.Shape
}

func (scene *Scene) Intersect(ray *math.Ray, thit *float32) bool {
	hit := false

	for _, shape := range scene.Shapes {
		if shape.Intersect(ray, thit) {
			ray.MaxT = *thit
			hit = true
		}
	}

	return hit
}