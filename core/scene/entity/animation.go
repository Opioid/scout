package entity

import (
	"github.com/Opioid/scout/base/math"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	_ "fmt"
)

type Animation struct {
	keyframes []keyframe
}

func (a *Animation) empty() bool {
	return len(a.keyframes) < 2
}

func (a *Animation) at(time float32) math.ComposedTransformation {
	f0 := &a.keyframes[0]
	f1 := &a.keyframes[1]

//	d := f1.time - f0.time

	fi := f0.transformation.Lerp(&f1.transformation, time)

	t := math.ComposedTransformation{}
	t.SetFromTransformation(&fi)
	return t
	

//	return MakeComposedTransformation(&fi)
} 

type keyframe struct {
	time float32
	transformation math.Transformation
}

func MakeAnimationFromJson(i interface{}) Animation {
	a := Animation{}

	keyframes, ok := i.([]interface{})

	if !ok {
		return a
	}

	a.keyframes = make([]keyframe, len(keyframes))

	for i, keyframe := range keyframes {
		a.keyframes[i] = parseKeyframe(keyframe)
	}

	return a
}

func parseKeyframe(i interface{}) keyframe {
	k := keyframe{}

	k.transformation.Scale = math.MakeIdentityVector3()
	k.transformation.Rotation = math.MakeIdentityQuaternion()

	keyframeNode, ok := i.(map[string]interface{})

	if !ok {
		return k
	}	

	for key, value := range keyframeNode {
		switch key {
		case "time":
			k.time = pkgjson.ParseFloat32(value)	
		case "position":
			k.transformation.Position = pkgjson.ParseVector3(value)
		case "scale":
			k.transformation.Scale = pkgjson.ParseVector3(value)
		case "rotation":
			k.transformation.Rotation = pkgjson.ParseRotationQuaternion(value)
		}
	}

	return k
}