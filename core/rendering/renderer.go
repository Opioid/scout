package rendering

import (
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/camera"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Renderer struct {

}

func (r *Renderer) Render(scene *pkgscene.Scene, context *Context) {
	target := context.Target
	dimensions := target.Dimensions()

	var ray math.Ray
	var intersection pkgscene.Intersection

	for y := 0; y < dimensions.Y; y++ {
		for x := 0; x < dimensions.X; x++ {
			context.Camera.GenerateRay(camera.NewSample(float32(x), float32(y)), &ray)

			if scene.Intersect(&ray, &intersection) {
				color := r.shade(scene, &intersection)
				target.Set(x, y, math.Vector4{color.X, color.Y, color.Z, 1.0})
			} else {
				target.Set(x, y, math.Vector4{0.0, 0.0, 0.0, 1.0})
			}
		}
	}


}

func (r *Renderer) shade(scene *pkgscene.Scene, intersection *pkgscene.Intersection) math.Vector3 {
	result := math.Vector3{0.0, 0.0, 0.0}

	color := intersection.Prop.Material.Color

	ray := math.Ray{Origin: intersection.Dg.Point, MaxT: 1000.0}

	for _, l := range scene.Lights {
		ray.Direction = l.Entity.Transformation.Matrix.Row(2).Vector3().Scale(-1.0)

		if !scene.IntersectP(&ray) {
			result = result.Add(color.Mul(l.Color))
		}
	}

	return result
}