package rendering

import (
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/camera"
	"github.com/Opioid/scout/base/math"
	"sync"
	_ "fmt"
)

type Renderer struct {
	BounceDepth int

	samplerDimensions math.Vector2i
	currentPixel math.Vector2i
}

func (r *Renderer) Render(scene *pkgscene.Scene, context *Context) {
	dimensions := context.Camera.Film().Dimensions()

	r.currentPixel = math.Vector2i{0, 0}

	r.samplerDimensions = math.Vector2i{16, 16}

	wg := sync.WaitGroup{}

	for {
		sampler := r.NewSubSampler(context.Sampler, dimensions)

		if sampler == nil {
			break
		}

		wg.Add(1)

		go func () {
			r.render(scene, context.Camera, sampler)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (r *Renderer) render(scene *pkgscene.Scene, camera camera.Camera, sampler pkgsampler.Sampler) {
	film := camera.Film()

	var ray math.OptimizedRay
	var sample pkgsampler.Sample

	for sampler.GenerateNewSample(&sample) {
		camera.GenerateRay(&sample, &ray)

		color := r.li(scene, &ray, 0) 

		film.AddSample(&sample, color)
	}
}

func (r *Renderer) li(scene *pkgscene.Scene, ray *math.OptimizedRay, depth int) math.Vector3 {
	var intersection prop.Intersection

	if scene.Intersect(ray, &intersection) {
		return r.shade(scene, ray, &intersection, depth)
	} else {
		return math.Vector3{0.0, 0.0, 0.0}
	}
}

func (r *Renderer) shade(scene *pkgscene.Scene, ray *math.OptimizedRay, intersection *prop.Intersection, depth int) math.Vector3 {
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

	if material.IsMirror() && depth < r.BounceDepth {
		reflection := intersection.Dg.Nn.Reflect(ray.Direction)

		secondaryRay := math.MakeOptimizedRay(intersection.Dg.P, reflection, intersection.Epsilon, 1000.0)

		result = result.Add(r.li(scene, &secondaryRay, depth + 1))
	}

	return result
}

func (r *Renderer) NewSubSampler(s pkgsampler.Sampler, dimensions math.Vector2i) pkgsampler.Sampler {
	if r.currentPixel.X >= dimensions.X {
		r.currentPixel.X = 0
		r.currentPixel.Y += r.samplerDimensions.Y
	}

	if r.currentPixel.Y >= dimensions.Y {
		return nil
	}

	end := r.currentPixel.Add(r.samplerDimensions).Min(dimensions)

	sampler := s.SubSampler(r.currentPixel, end)

	r.currentPixel.X += r.samplerDimensions.X

	return sampler
}