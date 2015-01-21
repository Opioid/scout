package glass

import (
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/base/math"
	_ "github.com/Opioid/math32"
	_ "math"
	_ "fmt"
)

type base struct {
	pool *Pool
}

func (b *base) Free(sample material.Sample, workerId uint32) {
	b.pool.Put(sample, workerId)
}

func (b *base) Energy() math.Vector3 {
	return math.MakeVector3(0.0, 0.0, 0.0)
}

func (b *base) IsLight() bool {
	return false
}

type Sample struct {
	values material.Values
}

func (s *Sample) Evaluate(l math.Vector3) math.Vector3 {
	return math.MakeVector3(0.3, 0.3, 0.3)
}

func (s *Sample) Values() material.Values {
	return s.values
}

