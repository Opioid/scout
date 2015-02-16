package glass

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
	_ "github.com/Opioid/math32"
	_ "math"
	_ "fmt"
)

type base struct {
	stack *BinnedStack
}

func (b *base) Free(sample material.Sample, workerID uint32) {
	b.stack.Push(workerID)
}

func (b *base) Energy() math.Vector3 {
	return math.MakeVector3(0.0, 0.0, 0.0)
}

func (b *base) IsLight() bool {
	return false
}

type Sample struct {
	material.SampleBase
	values material.Values
}

func NewSample() *Sample {
	s := &Sample{}
//	s.lambert.sample = s
//	s.ggx.sample = s
	return s
}

func (s *Sample) Evaluate(l math.Vector3) math.Vector3 {
	return math.MakeVector3(0.3, 0.3, 0.3)
}

func (s *Sample) Values() *material.Values {
	return &s.values
}

func (s *Sample) MonteCarloBxdf(subsample uint32, sampler sampler.Sampler) (material.Bxdf, float32) {
	return nil, 1.0
}