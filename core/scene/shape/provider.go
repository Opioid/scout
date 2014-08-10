package shape

import (
	"github.com/Opioid/scout/base/math"
	"io/ioutil"
	"encoding/json"
)

type Provider struct {

}

func (p *Provider) Load(filename string) Shape {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil
	}

	var document interface{}
	if err = json.Unmarshal(data, &document); err != nil {
		return nil
	}

	root := document.(map[string]interface{})

	t, hasType := root["type"]

	if !hasType {
		return nil
	}

	switch t {
	case "Sphere":
		return new(sphere)
	case "Plane":
		return new(plane)
	case "Triangle_mesh":
		return loadTriangleMesh(t)
	default:
		return nil
	}

	return nil
}

func loadTriangleMesh(i interface{}) *triangleMesh {
//	propNode, ok := i.(map[string]interface{})

	m := NewTriangleMesh(6, 4)

	m.setIndex(0, 0)
	m.setIndex(1, 1)
	m.setIndex(2, 2)

	m.setIndex(3, 2)
	m.setIndex(4, 1)
	m.setIndex(5, 3)

	m.setPosition(0, math.Vector3{-0.4, 1.8, 2.0})
	m.setNormal(0, math.Vector3{0.0, 0.0, -1.0})
	m.setUV(0, math.Vector2{0.0, 0.0})

	m.setPosition(1, math.Vector3{ 0.4, 1.8, 2.0})
	m.setNormal(1, math.Vector3{0.0, 0.0, -1.0})
	m.setUV(1, math.Vector2{1.0, 0.0})

	m.setPosition(2, math.Vector3{-0.4, 1.0, 2.0})
	m.setNormal(2, math.Vector3{0.0, 0.0, -1.0})
	m.setUV(2, math.Vector2{0.0, 1.0})

	m.setPosition(3, math.Vector3{ 0.4, 1.0, 2.0})
	m.setNormal(3, math.Vector3{0.0, 0.0, -1.0})
	m.setUV(3, math.Vector2{1.0, 1.0})

	return m
}
