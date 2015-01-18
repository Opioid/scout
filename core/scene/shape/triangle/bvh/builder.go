package bvh

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/shape/triangle/primitive"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	gomath "math"
	_ "fmt"
)

type Builder struct {

}

func (b *Builder) Build(triangles []primitive.IndexTriangle, vertices []geometry.Vertex, maxPrimitives int) Tree {
	primitiveIndices := make([]uint32, len(triangles))

	for i := range primitiveIndices {
		primitiveIndices[i] = uint32(i)
	}

	outTriangles := make([]primitive.Triangle, 0, len(triangles))

	root := buildNode{}
	root.split(primitiveIndices, triangles, vertices, maxPrimitives, 0, &outTriangles)

	return Tree{root, outTriangles}
}

type buildNode struct {
	aabb bounding.AABB

	startIndex, endIndex uint32

	children [2]*buildNode

	axis int8
}

func (n *buildNode) split(primitiveIndices []uint32, triangles []primitive.IndexTriangle, vertices []geometry.Vertex, maxPrimitives, depth int, 
						   outTriangles *[]primitive.Triangle) {
	n.aabb = subMeshAabb(primitiveIndices, triangles, vertices)

	if len(primitiveIndices) < maxPrimitives || depth > 18 {
		n.assign(primitiveIndices, triangles, vertices, outTriangles)
	} else {
		sp, axis := chooseSplittingPlane(&n.aabb)

		n.axis = axis

		numPids := len(primitiveIndices) / 2
		pids0 := make([]uint32, 0, numPids)
		pids1 := make([]uint32, 0, numPids)

		for _, pi := range primitiveIndices {
			s := triangleSide(vertices[triangles[pi].A].P, vertices[triangles[pi].B].P, vertices[triangles[pi].C].P, sp)
			
			if s == 0 {
				pids0 = append(pids0, pi)
			} else {
				pids1 = append(pids1, pi)
			}
		}

		n.children[0] = new(buildNode)
		n.children[1] = new(buildNode)

		n.children[0].split(pids0, triangles, vertices, maxPrimitives, depth + 1, outTriangles)

		pids0 = nil

		n.children[1].split(pids1, triangles, vertices, maxPrimitives, depth + 1, outTriangles)
	}
}

func (n *buildNode) assign(primitiveIndices []uint32, triangles []primitive.IndexTriangle, vertices []geometry.Vertex, 
							outTriangles *[]primitive.Triangle) {
	n.startIndex = uint32(len(*outTriangles))

	for _, pi := range primitiveIndices {
		*outTriangles = append(*outTriangles, primitive.MakeTriangle(&vertices[triangles[pi].A], 
											   						 &vertices[triangles[pi].B], 
																	 &vertices[triangles[pi].C], 
																	 triangles[pi].MaterialId))
	}

	n.endIndex = uint32(len(*outTriangles))
}

func (n *buildNode) intersect(ray *math.OptimizedRay, triangles []primitive.Triangle, intersection *primitive.Intersection) bool {
	if !n.aabb.IntersectP(ray) {
		return false
	}

	hit := false

	if n.children[0] != nil {
		c := ray.Sign[n.axis]

		if n.children[c].intersect(ray, triangles, intersection) {
			hit = true
		} 

		if n.children[1 - c].intersect(ray, triangles, intersection) {
			hit = true
		}
	} else {
		ti := primitive.Intersection{}
		ti.T = ray.MaxT

		for i := n.startIndex; i < n.endIndex; i++ {
			if h, c := triangles[i].Intersect(ray); h && c.T < ti.T {
				ti.Coordinates = c
				ti.Index = i
				hit = true
			}
		}

		if hit {
			// the idea was not to write these pointers in the loop... Don't know whether it makes a difference
			*intersection = ti
			ray.MaxT = ti.T
		}
	}

	return hit
}

func (n *buildNode) intersectP(ray *math.OptimizedRay, triangles []primitive.Triangle) bool {
	if !n.aabb.IntersectP(ray) {
		return false
	}

	if n.children[0] != nil {
		c := ray.Sign[n.axis]

		if n.children[c].intersectP(ray, triangles) {
			return true
		} 

		return n.children[1 - c].intersectP(ray, triangles)
	}

	for i := n.startIndex; i < n.endIndex; i++ {
		if triangles[i].IntersectP(ray) {
			return true
		}
	}

	return false
}

const (
	epsilon = 0.000000001
)

func subMeshAabb(primitiveIndices[]uint32, triangles []primitive.IndexTriangle, vertices []geometry.Vertex) bounding.AABB {
	min := math.MakeVector3( gomath.MaxFloat32,  gomath.MaxFloat32,  gomath.MaxFloat32)
	max := math.MakeVector3(-gomath.MaxFloat32, -gomath.MaxFloat32, -gomath.MaxFloat32)
	
	for _, pi := range primitiveIndices {
		min = triangleMin(vertices[triangles[pi].A].P, vertices[triangles[pi].B].P, vertices[triangles[pi].C].P, min)
		max = triangleMax(vertices[triangles[pi].A].P, vertices[triangles[pi].B].P, vertices[triangles[pi].C].P, max)
	}

	max.X += epsilon
	max.Y += epsilon
	max.Z += epsilon

	return bounding.MakeAABB(min, max)
}

func triangleMin(a, b, c, x math.Vector3) math.Vector3 {
	return a.Min(b).Min(c).Min(x)
}

func triangleMax(a, b, c, x math.Vector3) math.Vector3 {
	return a.Max(b).Max(c).Max(x)
}

var axis = [...]math.Vector3{ 
	math.MakeVector3(1.0, 0.0, 0.0),
	math.MakeVector3(0.0, 1.0, 0.0),
	math.MakeVector3(0.0, 0.0, 1.0), 
}

func triangleSide(a, b, c math.Vector3, p math.Plane) int {
	behind := 0

	if p.Behind(a) {
		behind++
	}

	if p.Behind(b) {
		behind++
	}

	if p.Behind(c) {
		behind++
	}

	if behind == 3 {
		return 0
	} else if behind == 0 {
		return 1
	} else {
		return 2
	}
}

func splittingPlane(aabb *bounding.AABB) (int8, float32) {
	position := aabb.Position()
	halfsize := aabb.Halfsize()

	if halfsize.X >= halfsize.Y && halfsize.X >= halfsize.Z {
		p := math.MakePlane(axis[0], position)
		return 0, p.D
	} else if halfsize.Y >= halfsize.X && halfsize.Y >= halfsize.Z {
		p := math.MakePlane(axis[1], position)
		return 1, p.D
	} else {
		p := math.MakePlane(axis[2], position)
		return 2, p.D
	}
}

func chooseSplittingPlane(aabb *bounding.AABB) (math.Plane, int8) {
	position := aabb.Position()
	halfsize := aabb.Halfsize()

	if halfsize.X >= halfsize.Y && halfsize.X >= halfsize.Z {
		return math.MakePlane(math.MakeVector3(1.0, 0.0, 0.0), position), 0
	} else if halfsize.Y >= halfsize.X && halfsize.Y >= halfsize.Z {
		return math.MakePlane(math.MakeVector3(0.0, 1.0, 0.0), position), 1
	} else {
		return math.MakePlane(math.MakeVector3(0.0, 0.0, 1.0), position), 2
	}
}