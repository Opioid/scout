package shape

import (
	_ "github.com/Opioid/scout/base/math"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"io/ioutil"
	"encoding/json"
	_ "fmt"
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

	if g, ok := root["geometry"]; ok {
		return loadGeometry(g)
	}

	return nil
}

func loadGeometry(i interface{}) Shape {
	geometryNode, ok := i.(map[string]interface{})

	if !ok {
		return nil
	}

	if geometryNode["primitive_topology"] != "triangle_list" {
		return nil
	}

	/*
	if g, ok := geometryNode["groups"]; ok {
		if groups, ok := g.([]interface{}); ok {

		}
	}

	grou
	*/

	var indices []interface{}
	var positions []interface{}

	if i, ok := geometryNode["indices"]; ok {
		indices = i.([]interface{})
	}

	if p, ok := geometryNode["positions"]; ok {
		positions = p.([]interface{})
	}

	if indices == nil || positions == nil {
		return nil
	}

	m := NewTriangleMesh(uint32(len(indices)), uint32(len(positions)))

	for i, index := range indices {
		m.setIndex(uint32(i), uint32(index.(float64)))
	}

	for i, position := range positions {
		m.setPosition(uint32(i), pkgjson.ParseVector3(position))
	}

	if n, ok := geometryNode["normals"]; ok {
		normals := n.([]interface{})

		for i, normal := range normals {
			m.setNormal(uint32(i), pkgjson.ParseVector3(normal))
		}
	}

	if u, ok := geometryNode["texture_coordinates_0"]; ok {
		uvs := u.([]interface{})

		for i, uv := range uvs {
			m.setUV(uint32(i), pkgjson.ParseVector2(uv))
		}
	}

	m.compile()

	return m
}