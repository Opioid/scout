package scene

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/core/scene/shape"
)

type Scene struct {
	Props []*Prop
	Shapes []shape.Shape
}

func (scene *Scene) Intersect(ray *math.Ray, thit *float32) bool {
	hit := false

	for _, prop := range scene.Props {
		if prop.Intersect(ray, thit) {
			ray.MaxT = *thit
			hit = true
		}
	}
/*
	transformation := &math.Transformation{}

	for _, shape := range scene.Shapes {
		if shape.Intersect(transformation, ray, thit) {
			ray.MaxT = *thit
			hit = true
		}
	}
*/
	return hit
}