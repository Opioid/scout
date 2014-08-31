package surrounding

import (
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
)

type sphere struct {
	sphereMap texture.SamplerSphere
}

func NewSphere(sphereMap texture.SamplerSphere) *sphere {
	return &sphere{sphereMap}
}

func (s *sphere) Sample(ray *math.OptimizedRay) math.Vector3 {
	return s.sphereMap.Sample(ray.Direction).Vector3()
}