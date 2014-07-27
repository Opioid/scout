package rendering

import (
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/camera"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Renderer struct {
	BounceDepth int
}

func (r *Renderer) Render(scene *pkgscene.Scene, context *Context) {
	target := context.Target
	dimensions := target.Dimensions()

	var ray math.Ray

	for y := 0; y < dimensions.Y; y++ {
		for x := 0; x < dimensions.X; x++ {
			context.Camera.GenerateRay(camera.NewSample(float32(x), float32(y)), &ray)

			color := r.render(scene, &ray, 0) 

			target.Set(x, y, math.Vector4{color.X, color.Y, color.Z, 1.0})
		}
	}
}

func (r *Renderer) render(scene *pkgscene.Scene, ray *math.Ray, depth int) math.Vector3 {
	var intersection pkgscene.Intersection

	if scene.Intersect(ray, &intersection) {
		return r.shade(scene, ray.Direction, &intersection, depth)
	} else {
		return math.Vector3{0.0, 0.0, 0.0}
	}
}

func (r *Renderer) shade(scene *pkgscene.Scene, eyeDirection math.Vector3, intersection *pkgscene.Intersection, depth int) math.Vector3 {
	result := math.Vector3{0.0, 0.0, 0.0}

	color := intersection.Prop.Material.Color
	roughness := intersection.Prop.Material.Roughness

	ray := math.Ray{Origin: intersection.Dg.P, MinT: intersection.Epsilon, MaxT: 1000.0}

	for _, l := range scene.Lights {
		if roughness > 0 {
			lightVector := l.Entity.Transformation.Matrix.Row(2).Vector3().Scale(-1.0)

			ray.Direction = lightVector

			if !scene.IntersectP(&ray) {
				cos := math.Max(intersection.Dg.Nn.Dot(lightVector), 0.0)

				diffuse := color.Mul(l.Color).Scale(cos)

				result = result.Add(diffuse)
			}
		} else if depth < r.BounceDepth {
			reflection := intersection.Dg.Nn.Reflect(eyeDirection)

			secondaryRay := math.Ray{Origin: intersection.Dg.P, Direction: reflection, MinT: intersection.Epsilon, MaxT: 1000.0}

			result = r.render(scene, &secondaryRay, depth + 1).Mul(color)
		}
	}

	return result
}