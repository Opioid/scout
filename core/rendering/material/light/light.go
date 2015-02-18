package light

import (
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type base struct {
	light light.Light
}

func (b *base) Free(sample material.Sample, workerID uint32) {
	// b.pool.Put(sample, workerID)
}

func (b *base) IsMirror() bool {
	return false
}

func (b *base) IsLight() bool {
	return true
}

type Sample struct {
	material.SampleBase
	values material.Values
}

func (s *Sample) Evaluate(l math.Vector3) math.Vector3 {
	return math.MakeVector3(0.0, 0.0, 0.0)
}

func (s *Sample) Values() *material.Values {
	return &s.values
}

func (s *Sample) MonteCarloBxdf(subsample uint32, sampler sampler.Sampler) (material.Bxdf, float32) {
	return nil, 1.0
}

func (s *Sample) SampleEvaluate(subsample uint32, sampler sampler.Sampler) (math.Vector3, math.Vector3, float32) {
	return math.MakeVector3(0.0, 0.0, 0.0), math.MakeVector3(0.0, 0.0, 0.0), 0.0
}