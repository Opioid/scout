package json

import (
	"github.com/Opioid/scout/base/math"
)

func ParseFloat32(value interface{}) float32 {
	return float32(value.(float64))
}

func ParseVector2(value interface{}) math.Vector2 {
	if floats, ok := value.([]interface{}); ok {
		return math.Vector2{float32(floats[0].(float64)), float32(floats[1].(float64))}
	} else {
		return math.Vector2{}
	}
}

func ParseVector2i(value interface{}) math.Vector2i {
	if ints, ok := value.([]interface{}); ok {
		return math.Vector2i{int(ints[0].(float64)), int(ints[1].(float64))}
	} else {
		return math.Vector2i{}
	}
}

func ReadVector2i(value map[string]interface{}, name string, defaultValue math.Vector2i) math.Vector2i {
	if t, ok := value[name]; ok {
		if ints, isArray := t.([]interface{}); isArray {
			return math.Vector2i{int(ints[0].(float64)), int(ints[1].(float64))}
		}
	}

	return defaultValue
}

func ParseVector3(value interface{}) math.Vector3 {
	if floats, ok := value.([]interface{}); ok {
		return math.Vector3{float32(floats[0].(float64)), float32(floats[1].(float64)), float32(floats[2].(float64))}
	} else {
		return math.Vector3{}
	}
}

func ParseRotationMatrix(value interface{}) math.Matrix3x3 {
	rotation := ParseVector3(value)

	rotationX := math.Matrix3x3{}
	rotationX.SetRotationX(math.DegreesToRadians(rotation.X))

	rotationY := math.Matrix3x3{}
	rotationY.SetRotationY(math.DegreesToRadians(rotation.Y))

	rotationZ := math.Matrix3x3{}
	rotationZ.SetRotationZ(math.DegreesToRadians(rotation.Z))

	t := rotationZ.Multiply(&rotationX)
	return t.Multiply(&rotationY)
}

func ParseRotationQuaternion(value interface{}) math.Quaternion {
	m := ParseRotationMatrix(value)
	return math.MakeQuaternionFromMatrix3x3(&m)
}