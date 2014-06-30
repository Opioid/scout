package json

import (
	"github.com/Opioid/scout/base/math"
)

func ParseVector2i(value interface{}) math.Vector2i {
	if ints, ok := value.([]interface{}); ok {
		return math.Vector2i{int(ints[0].(float64)), int(ints[1].(float64))}
	} else {
		return math.Vector2i{}
	}
}

func ParseVector3(value interface{}) math.Vector3 {
	if floats, ok := value.([]interface{}); ok {
		return math.Vector3{float32(floats[0].(float64)), float32(floats[1].(float64)), float32(floats[2].(float64))}
	} else {
		return math.Vector3{}
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