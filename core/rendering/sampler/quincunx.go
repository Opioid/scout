package sampler

import (
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Quincunx struct {
	sampler
	currentPixel math.Vector2i
	currentSample int
}

var quincunxOffsets = []math.Vector2{
	math.MakeVector2(0.25, 0.25),
	math.MakeVector2(0.75, 0.25),
	math.MakeVector2(0.5,  0.5),
	math.MakeVector2(0.25, 0.75),
	math.MakeVector2(0.75, 0.75),
}

func NewQuincunx(start, end math.Vector2i) *Quincunx {
	q := new(Quincunx)
	q.Resize(start, end)
	q.currentPixel = start
	return q
}

func (q *Quincunx) Resize(start, end math.Vector2i) {
	q.start = start
	q.end = end
}

func (q *Quincunx) Restart() {
	q.currentPixel = q.start
	q.currentSample = 0
}

func (q *Quincunx) SubSampler(start, end math.Vector2i) Sampler {
	return NewQuincunx(start, end)
}

func (q *Quincunx) GenerateNewSample(s *Sample) bool {
	if q.currentPixel.X >= q.end.X {
		q.currentPixel.X = q.start.X
		q.currentPixel.Y++
	}

	if q.currentPixel.Y >= q.end.Y {
		return false
	}

	o := quincunxOffsets[q.currentSample]

	s.Coordinates = math.MakeVector2(float32(q.currentPixel.X) + o.X, float32(q.currentPixel.Y) + o.Y)
	s.Id = uint32(q.currentSample)

	q.currentSample++

	if q.currentSample >= 5 {
		q.currentSample = 0
		q.currentPixel.X++
	}

	return true
}

func (q *Quincunx) NumSamplesPerPixel() uint32 {
	return 5
}