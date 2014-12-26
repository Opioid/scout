package shape

import (
	"github.com/Opioid/scout/core/scene/shape/triangle"
	_ "github.com/Opioid/scout/base/math"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"io/ioutil"
	"encoding/json"
	_ "fmt"
)

type Provider struct {
	meshes map[string]Shape
}

func NewProvider() *Provider {
	p := Provider{}
	p.meshes = make(map[string]Shape)
	return &p
}

func (p *Provider) Load(filename string) Shape {
	if mesh, ok := p.meshes[filename]; ok {
		return mesh
	}

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil
	}

	var document interface{}
	if err = json.Unmarshal(data, &document); err != nil {
		return nil
	}

	root := document.(map[string]interface{})

	var mesh Shape

	if g, ok := root["geometry"]; ok {
		mesh = loadGeometry(g)
	}

	if mesh != nil {
		p.meshes[filename] = mesh
	}

	return mesh
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

	m := triangle.NewMesh(uint32(len(indices)), uint32(len(positions)))

	for i, index := range indices {
		m.SetIndex(uint32(i), uint32(index.(float64)))
	}

	for i, position := range positions {
		m.SetPosition(uint32(i), pkgjson.ParseVector3(position))
	}

	if n, ok := geometryNode["normals"]; ok {
		normals := n.([]interface{})

		for i, normal := range normals {
			m.SetNormal(uint32(i), pkgjson.ParseVector3(normal))
		}
	}

	if t, ok := geometryNode["tangents_and_bitangent_signs"]; ok {
		tangents := t.([]interface{})

		for i, tangent := range tangents {
			tas := pkgjson.ParseVector4(tangent)
			m.SetTangentAndSign(uint32(i), tas.Vector3(), tas.W)
		}
	}

	if u, ok := geometryNode["texture_coordinates_0"]; ok {
		uvs := u.([]interface{})

		for i, uv := range uvs {
			m.SetUV(uint32(i), pkgjson.ParseVector2(uv))
		}
	}

	m.Compile()

	return m
}