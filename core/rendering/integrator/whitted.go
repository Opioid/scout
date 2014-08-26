package integrator

import (
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
)

type whitted struct {
	bounceDepth int
}

func NewWhitted(bounceDepth int) *whitted {
	return &whitted{bounceDepth}
}

func (w *whitted) Li(scene *pkgscene.Scene, ray *math.OptimizedRay, rng *random.Generator) math.Vector3 {
	return w.render(scene, ray, 0)
}

func (w *whitted) render(scene *pkgscene.Scene, ray *math.OptimizedRay, depth int) math.Vector3 {
	var intersection prop.Intersection

	if scene.Intersect(ray, &intersection) {
		return w.shade(scene, ray, &intersection, depth)
	} else {
		return scene.Surrounding.Sample(ray)
	}
}

func (w *whitted) shade(scene *pkgscene.Scene, ray *math.OptimizedRay, intersection *prop.Intersection, depth int) math.Vector3 {
	result := math.Vector3{0.0, 0.0, 0.0}

	material := intersection.Prop.Material

	shadowRay := math.OptimizedRay{}
	shadowRay.Origin = intersection.Dg.P
	shadowRay.MinT = intersection.Epsilon
	shadowRay.MaxT = 1000.0

	v := ray.Direction.Scale(-1.0)

	for _, l := range scene.Lights {
		lightVector := l.Vector(intersection.Dg.P)

		shadowRay.SetDirection(lightVector)

		if !scene.IntersectP(&shadowRay) {
			color, opacity := material.Evaluate(&intersection.Dg, lightVector, v)

		//	result = result.Add(color.Scale(opacity).Mul(l.Color))

			result.AddAssign(l.Light(intersection.Dg.P, color.Scale(opacity)))

/*
			if opacity < 1.0 {
				secondaryRay := *ray
				secondaryRay.MinT = ray.MaxT + intersection.Epsilon
				secondaryRay.MaxT = 1000.0
				
				secondaryColor := r.li(scene, &secondaryRay, depth)
	
				result = result.Add(secondaryColor.Scale(1.0 - opacity))
			}
			*/
		}
	}

	if material.IsMirror() && depth < w.bounceDepth {
		reflection := intersection.Dg.Nn.Reflect(ray.Direction)

		secondaryRay := math.MakeOptimizedRay(intersection.Dg.P, reflection, intersection.Epsilon, 1000.0)

		result = result.Add(w.render(scene, &secondaryRay, depth + 1))
	}

	return result
}