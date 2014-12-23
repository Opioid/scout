package entity

import (
	"github.com/Opioid/scout/base/math"
)

type animation struct {
	keyframes []keyframe
}

func (a *animation) empty() bool {
	return 0 == len(a.keyframes)
}

func (a *animation) at(time float32, t *ComposedTransformation) {
	t.SetFromTransformation(&a.keyframes[0].transformation)
} 

type keyframe struct {
	time float32
	transformation math.Transformation
}