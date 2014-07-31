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
	currentSampler, numSamplers math.Vector2i
}

func (r *Renderer) Render(scene *pkgscene.Scene, context *Context) {
	dimensions := context.Camera.Film().Dimensions()

	r.samplerDimensions = math.Vector2i{16, 16}
	r.numSamplers = math.Vector2i{dimensions.X / r.samplerDimensions.X, dimensions.Y / r.samplerDimensions.Y}

	wg := sync.WaitGroup{}

	for {
		sampler := r.NewSubSampler(context.Sampler)

		if sampler == nil {
			break
		}

		wg.Add(1)
		go r.render(scene, context.Camera, sampler, &wg)
	}

	wg.Wait()
}

func (r *Renderer) render(scene *pkgscene.Scene, camera camera.Camera, sampler pkgsampler.Sampler, wg *sync.WaitGroup) {
	film := camera.Film()

	var ray math.Ray
	var sample pkgsampler.Sample

	for sampler.GenerateNewSample(&sample) {
		camera.GenerateRay(&sample, &ray)

		color := r.li(scene, &ray, 0) 

		film.AddSample(&sample, color)
	}

	wg.Done()
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

			result = r.li(scene, &secondaryRay, depth + 1).Mul(color)
		}
	}

	return result
}

func (r *Renderer) NewSubSampler(s pkgsampler.Sampler) pkgsampler.Sampler {
	if r.currentSampler.X >= r.numSamplers.X {
		r.currentSampler.X = 0
		r.currentSampler.Y++
	}

	if r.currentSampler.Y >= r.numSamplers.Y {
		return nil
	}

	start := math.Vector2i{r.currentSampler.X * r.samplerDimensions.X, r.currentSampler.Y * r.samplerDimensions.Y}

	sampler := s.SubSampler(start, start.Add(r.samplerDimensions))

	r.currentSampler.X++

	return sampler
}