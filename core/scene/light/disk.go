package light

import (
	"github.com/Opioid/scout/core/rendering/sampler"
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

func (l *Disk) Samples(p math.Vector3, time float32, subsample, maxSamples uint32, sampler *sampler.ScrambledHammersley, samples *[]Sample) {
	transformation := l.entity.TransformationAt(time)

	result := Sample{}

	tsamples := sampler.GenerateSamples(subsample)

	for _, sample := range tsamples {
		ls := math.DiskSample_uniform(sample.X, sample.Y)
		ws := transformation.Rotation.TransformVector3(ls).Scale(l.radius)

		v := transformation.Rotation.Direction().Scale(-1.0).Add(ws)

		result.L = v.Normalized()
		result.Energy = l.color

		*samples = append(*samples, result)
	}	
}
