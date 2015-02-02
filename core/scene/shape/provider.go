package shape

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/shape/triangle"
	"github.com/Opioid/scout/core/scene/shape/triangle/bvh"
	"github.com/Opioid/scout/core/scene/shape/triangle/primitive"
	"github.com/Opioid/scout/base/math"
	_ "github.com/Opioid/scout/base/parsing/json"
	"io/ioutil"
	_ "os"
	"encoding/json"
	"runtime"
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

	triangles, vertices := loadMeshData(filename)

	runtime.GC()

	if triangles == nil || vertices == nil {
		return nil
	}

	builder := bvh.Builder{}
	tree := builder.Build(triangles, vertices, 8)

	mesh := triangle.NewMesh(tree.AABB(), tree)	

	p.meshes[filename] = mesh

	return mesh
}

type jsonMesh struct {
	Geometry Geometry
}

type Geometry struct {
	Groups []group

	Primitive_topology string

	Indices []uint32

	Positions [][3]float32
	Normals   [][3]float32
	Tangents_and_bitangent_signs [][4]float32

	Texture_coordinates_0 [][2]float32
}

type group struct {
	Material_index uint32
	Start_index    uint32
	Num_indices    uint32
}

func loadMeshData(filename string) ([]primitive.IndexTriangle, []geometry.Vertex) {
	
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil
	}

	var mesh jsonMesh
	if err = json.Unmarshal(data, &mesh); err != nil {
		return nil, nil
	}	
	
	/*
	fi, err := os.Open(filename)

	defer fi.Close()

	if err != nil {
		return nil, nil
	}	

	var mesh jsonMesh
	if err = json.NewDecoder(fi).Decode(&mesh); err != nil {
		return nil, nil
	}
	*/

	if mesh.Geometry.Primitive_topology != "triangle_list" {
		return nil, nil
	}

	numTriangles := uint32(len(mesh.Geometry.Indices)) / 3

	triangles := make([]primitive.IndexTriangle, 0, numTriangles)

	maxMaterialId := uint32(len(mesh.Geometry.Groups) - 1)

	for _, p := range mesh.Geometry.Groups {
		trianglesStart := p.Start_index / 3
		trianglesEnd := (p.Start_index + p.Num_indices) / 3

		for i := trianglesStart; i < trianglesEnd; i++ {
			a := mesh.Geometry.Indices[i * 3 + 0]
			b := mesh.Geometry.Indices[i * 3 + 1]
			c := mesh.Geometry.Indices[i * 3 + 2]

			triangles = append(triangles, primitive.MakeIndexTriangle(a, b, c, math.Minui(p.Material_index, maxMaterialId)))
		}
	}

	numVertices := len(mesh.Geometry.Positions)

	if numVertices == 0 {
		return nil, nil
	}

	vertices := make([]geometry.Vertex, numVertices)

	for i := range mesh.Geometry.Positions {
		vertices[i].P = math.MakeVector3(mesh.Geometry.Positions[i][0], mesh.Geometry.Positions[i][1], mesh.Geometry.Positions[i][2])
	}

	for i := range mesh.Geometry.Normals {
		vertices[i].N = math.MakeVector3(mesh.Geometry.Normals[i][0], mesh.Geometry.Normals[i][1], mesh.Geometry.Normals[i][2])
	}

	for i := range mesh.Geometry.Tangents_and_bitangent_signs {
		vertices[i].T = math.MakeVector3(mesh.Geometry.Tangents_and_bitangent_signs[i][0], 
										 mesh.Geometry.Tangents_and_bitangent_signs[i][1], 
										 mesh.Geometry.Tangents_and_bitangent_signs[i][2])

		vertices[i].BitangentSign = mesh.Geometry.Tangents_and_bitangent_signs[i][3]
	}

	for i := range mesh.Geometry.Texture_coordinates_0 {
		vertices[i].UV = math.MakeVector2(mesh.Geometry.Texture_coordinates_0[i][0], mesh.Geometry.Texture_coordinates_0[i][1])
	}

	if len(mesh.Geometry.Normals) > 0 && len(mesh.Geometry.Tangents_and_bitangent_signs) == 0 {
		// If normals but no tangents were loaded, compute the tangent space manually

		basis := math.Matrix3x3{}

		for i := range vertices {
			basis.SetBasis(vertices[i].N)

			vertices[i].T = basis.Right()
			vertices[i].BitangentSign = 1.0
		}
	}

	return triangles, vertices
}