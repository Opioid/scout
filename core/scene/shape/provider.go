package shape

import (
	"github.com/Opioid/scout/core/scene/shape/triangle"
	"github.com/Opioid/scout/core/scene/shape/triangle/primitive"
	"github.com/Opioid/scout/base/math"
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
	
	/*
	fi, err := os.Open(filename)

	defer fi.Close()

	if err != nil {
		return nil
	}	

	var document interface{}
	if err = json.NewDecoder(fi).Decode(&document); err != nil {
		return nil
	}
	*/
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

	var groups []interface{}
	var indices []interface{}
	var positions []interface{}

	if g, ok := geometryNode["groups"]; ok {
		groups = g.([]interface{})
	}

	if i, ok := geometryNode["indices"]; ok {
		indices = i.([]interface{})
	}

	if p, ok := geometryNode["positions"]; ok {
		positions = p.([]interface{})
	}

	if groups == nil || indices == nil || positions == nil {
		return nil
	}

	numTriangles := uint32(len(indices)) / 3

	m := triangle.NewMesh(numTriangles, uint32(len(positions)))

	maxMaterialId := uint32(len(groups) - 1)

	for _, g := range groups {
		if groupNode, ok := g.(map[string]interface{}); ok {
			start    := pkgjson.ReadUint32(groupNode, "start_index", 0)
			count    := pkgjson.ReadUint32(groupNode, "num_indices", 0)
			material := pkgjson.ReadUint32(groupNode, "material_index", 0)

			trianglesStart := start / 3
			trianglesEnd := (start + count) / 3

			for i := trianglesStart; i < trianglesEnd; i++ {
				a := uint32(indices[i * 3 + 0].(float64))
				b := uint32(indices[i * 3 + 1].(float64))
				c := uint32(indices[i * 3 + 2].(float64))
				m.SetTriangle(i, primitive.MakeIndexTriangle(a, b, c, math.Minui(material, maxMaterialId)))
			}
		}
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