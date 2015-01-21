package json

import (
	"github.com/Opioid/scout/base/math"
)

func ReadUint32(value map[string]interface{}, name string, defaultValue uint32) uint32 {
	if t, ok := value[name]; ok {
		if f, isFloat := t.(float64); isFloat {
			return uint32(f)
		}
	}

	return defaultValue
}

func ParseFloat32(value interface{}) float32 {
	return float32(value.(float64))
}

func ReadFloat32(value map[string]interface{}, name string, defaultValue float32) float32 {
	if t, ok := value[name]; ok {
		if f, isFloat := t.(float64); isFloat {
			return float32(f)
		}
	}

	return defaultValue
}

func ParseVector2(value interface{}) math.Vector2 {
	if floats, ok := value.([]interface{}); ok {
		return math.MakeVector2(float32(floats[0].(float64)), float32(floats[1].(float64)))
	} 
		
	return math.Vector2{}
}

func ParseVector2i(value interface{}) math.Vector2i {
	if ints, ok := value.([]interface{}); ok {
		return math.MakeVector2i(int32(ints[0].(float64)), int32(ints[1].(float64)))
	} 
		
	return math.Vector2i{}
}

func ReadVector2i(value map[string]interface{}, name string, defaultValue math.Vector2i) math.Vector2i {
	if t, ok := value[name]; ok {
		if ints, isArray := t.([]interface{}); isArray {
			return math.MakeVector2i(int32(ints[0].(float64)), int32(ints[1].(float64)))
		}
	}

	return defaultValue
}

func ParseVector3(value interface{}) math.Vector3 {
	if floats, ok := value.([]interface{}); ok {
		return math.MakeVector3(float32(floats[0].(float64)), float32(floats[1].(float64)), float32(floats[2].(float64)))
	} 
		
	return math.Vector3{}
}

func ReadVector3(value map[string]interface{}, name string, defaultValue math.Vector3) math.Vector3 {
	if t, ok := value[name]; ok {
		if floats, isArray := t.([]interface{}); isArray {
			return math.MakeVector3(float32(floats[0].(float64)), float32(floats[1].(float64)), float32(floats[2].(float64)))
		}
	}

	return defaultValue
}

func ParseVector4(value interface{}) math.Vector4 {
	if floats, ok := value.([]interface{}); ok {
		return math.MakeVector4(float32(floats[0].(float64)), float32(floats[1].(float64)), float32(floats[2].(float64)), float32(floats[3].(float64)))
	} 
		
	return math.Vector4{}
}

func ReadString(value map[string]interface{}, name, defaultValue string) string {
	if t, ok := value[name]; ok {
		if s, isString := t.(string); isString {
			return s
		}
	}

	return defaultValue
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

func MakeRotationMatrix(a [3]float32) math.Matrix3x3 {
	rotation := math.MakeVector3FromArray(a)

	rotationX := math.Matrix3x3{}
	rotationX.SetRotationX(math.DegreesToRadians(rotation.X))

	rotationY := math.Matrix3x3{}
	rotationY.SetRotationY(math.DegreesToRadians(rotation.Y))

	rotationZ := math.Matrix3x3{}
	rotationZ.SetRotationZ(math.DegreesToRadians(rotation.Z))

	t := rotationZ.Multiply(&rotationX)
	return t.Multiply(&rotationY)
}

func MakeRotationQuaternion(a [3]float32) math.Quaternion {
	m := MakeRotationMatrix(a)
	return math.MakeQuaternionFromMatrix3x3(&m)
}