package integrator

import (
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	gomath "math"
	_ "fmt"
)

type ao struct {
	numSamples int
	radius float32
}

func NewAo(numSamples int, radius float32) *ao {
	return &ao{numSamples, radius}
}

func (a *ao) Li(scene *pkgscene.Scene, ray *math.OptimizedRay, rng *random.Generator) math.Vector3 {
	return a.render(scene, ray, rng)
}

func (a *ao) render(scene *pkgscene.Scene, ray *math.OptimizedRay, rng *random.Generator) math.Vector3 {
	var intersection prop.Intersection

	if scene.Intersect(ray, &intersection) {
		return a.shade(scene, ray, rng, &intersection)
	} else {
		return scene.Surrounding.Sample(ray)
	}
}

func (a *ao) shade(scene *pkgscene.Scene, ray *math.OptimizedRay, rng *random.Generator, intersection *prop.Intersection) math.Vector3 {
	occlusionRay := math.OptimizedRay{}
	occlusionRay.Origin = intersection.Dg.P
	occlusionRay.MinT = intersection.Epsilon
	occlusionRay.MaxT = a.radius

	result := float32(0.0)

	for i := 0; i < a.numSamples; i++ {
		v := hemisphereSample(intersection.Dg.Nn, rng.RandomFloat32(), rng.RandomFloat32())
		occlusionRay.SetDirection(v)

		if !scene.IntersectP(&occlusionRay) {
			result += 1.0
		}
	}

	result = result / float32(a.numSamples)

	return math.Vector3{result, result, result}
}

func hemisphereSample(v math.Vector3, r, s float32) math.Vector3 {
	basis := math.Matrix3x3{}
	basis.SetBasis(v)

	r1 := s * 2.0 * gomath.Pi
	r2 := math.Sqrt(1.0 - r)
	sp := math.Sqrt(1.0 - r2 * r2)
	dir := math.Vector3{math.Cos(r1) * sp, math.Sin(r1) * sp, r2}

//	fmt.Printf("%f %f\n", r, s)

	return basis.TransformVector3(dir)
}

/*
Vector3 Raytracer::getHemisphereSample(const Vector3 &v, float r, float s)
{
	Matrix3 basis;
	setBasis(basis, v);

	Vector3 dir;
	float r1 = s * 2.f * math::pi;
	float r2 = sqrt(1.f - r);
	float sp = sqrt(1.f - r2 * r2); // !! == sqrtf(u) !!
	dir.x = cos(r1) * sp;
	dir.y = sin(r1) * sp;
	dir.z = r2;

//	cosTheta = r2;
//	pdf = cosTheta*(1.0f/PI_f);

	return dir * basis;
}
*/