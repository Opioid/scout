package sampler

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	_ "fmt"
)

type Quincunx struct {
	currentSample uint32
}

var quincunxSamples = []math.Vector2{
	math.MakeVector2(0.25, 0.25),
	math.MakeVector2(0.75, 0.25),
	math.MakeVector2(0.5,  0.5),
	math.MakeVector2(0.25, 0.75),
	math.MakeVector2(0.75, 0.75),
}

func NewQuincunx() *Quincunx {
	q := new(Quincunx)
	return q
}

func (q *Quincunx) Clone(rng *random.Generator) Sampler {
	return NewQuincunx()
}

func (q *Quincunx) NumSamplesPerIteration() uint32 {
	return 5
}

func (q *Quincunx) Restart(numIterations uint32) {
	q.currentSample = 0
}

func (q *Quincunx) GenerateNewSample(offset math.Vector2, sample *CameraSample) bool {
	if q.currentSample >= 5 {
		return false
	}

	sample.Coordinates = quincunxSamples[q.currentSample]

	q.currentSample++

	return true
}

func (q *Quincunx) GenerateSamples(iteration uint32) []math.Vector2 {
	return quincunxSamples
}