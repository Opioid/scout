package light

import (
	pkgsampler "github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
)

type Disk struct {
	light
	radius float32
}

func NewDisk(radius float32) *Disk {
	d := Disk{}

	d.radius = radius

	return &d
}

func (l *Disk) Samples(p math.Vector3, sampler *pkgsampler.Stratified, samples *[]Sample) {
/*	s := Sample{}

	s.L = l.entity.Transformation.Rotation.Direction().Scale(-1.0)
	s.Energy = l.color

	*samples = append(*samples, s)
*/
	result := Sample{}

	sample := pkgsampler.Sample{}
//	for i, len := 0, cap(*samples); i < len; i++ {
	for sampler.GenerateNewSample(&sample) {
		ls := math.DiskSample_uniform(sample.Coordinates.X, sample.Coordinates.Y)
		ws := l.entity.Transformation.Rotation.TransformVector3(ls).Scale(l.radius)

		v := l.entity.Transformation.Rotation.Direction().Scale(-1.0).Add(ws)

		result.L = v.Normalized()
		result.Energy = l.color

		*samples = append(*samples, result)
	}	
}
