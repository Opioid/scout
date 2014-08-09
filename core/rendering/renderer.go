package rendering

import (
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	pkgscene "github.com/Opioid/scout/core/scene"
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

	var ray math.Ray
	var sample pkgsampler.Sample

	for sampler.GenerateNewSample(&sample) {
		camera.GenerateRay(&sample, &ray)

		color := r.li(scene, &ray, 0) 

		film.AddSample(&sample, color)
	}

	
}

func (r *Renderer) li(scene *pkgscene.Scene, ray *math.Ray, depth int) math.Vector3 {
	var intersection pkgscene.Intersection

	if scene.Intersect(ray, &intersection) {
		return r.shade(scene, ray.Direction, &intersection, depth)
	} else {
		return math.Vector3{0.0, 0.0, 0.0}
	}
}

func (r *Renderer) shade(scene *pkgscene.Scene, eyeDirection math.Vector3, intersection *pkgscene.Intersection, depth int) math.Vector3 {
	result := math.Vector3{0.0, 0.0, 0.0}

	material := intersection.Prop.Material

	ray := math.Ray{Origin: intersection.Dg.P, MinT: intersection.Epsilon, MaxT: 1000.0}

	v := eyeDirection.Scale(-1.0)

	for _, l := range scene.Lights {
		lightVector := l.Entity.Transformation.Matrix.Row(2).Vector3().Scale(-1.0)

		ray.Direction = lightVector

		if !scene.IntersectP(&ray) {
		//	diffuse := color.Mul(l.Color).Scale(cos)
			color := material.Evaluate(intersection.Dg.Nn, lightVector, v).Mul(l.Color)

			result = result.Add(color)
		}
	}

	if material.IsMirror() && depth < r.BounceDepth {
		reflection := intersection.Dg.Nn.Reflect(eyeDirection)

		secondaryRay := math.Ray{Origin: intersection.Dg.P, Direction: reflection, MinT: intersection.Epsilon, MaxT: 1000.0}

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