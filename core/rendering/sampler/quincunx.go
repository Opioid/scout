package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type Quincunx struct {
	currentPixel math.Vector2i
	currentSample uint32
}

var quincunxOffsets = []math.Vector2{
	math.MakeVector2(-0.5, -0.5),
	math.MakeVector2( 0.5, -0.5),
	math.MakeVector2( 0.0,  0.0),
	math.MakeVector2(-0.5,  0.5),
	math.MakeVector2( 0.5,  0.5),
}

func NewQuincunx() *Quincunx {
	q := new(Quincunx)
	return q
}

func (q *Quincunx) Clone(rng *random.Generator) Sampler {
	return NewQuincunx()
}

func (q *Quincunx) Restart() {
	q.currentSample = 0
}

func (q *Quincunx) GenerateNewSample(sample *math.Vector2) bool {
	if q.currentSample >= 5 {
		return false
	}

	*sample = quincunxOffsets[q.currentSample]

	q.currentSample++

	return true
}

func (q *Quincunx) NumSamplesPerPixel() uint32 {
	return 5
}