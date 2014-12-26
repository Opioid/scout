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

func (b *Builder) Build(indices []uint32, vertices []geometry.Vertex, maxPrimitives int, tree *Tree) {
	primitiveIndices := make([]uint32, len(indices) / 3)

	for i := range primitiveIndices {
		primitiveIndices[i] = uint32(i) * 3
	}

	root := buildNode{}
	root.split(primitiveIndices, indices, vertices, maxPrimitives, 0)

	tree.root = root
}

type buildNode struct {
	aabb bounding.AABB

	axis int

//	indices []uint32
	triangles []primitive.Triangle

	children [2]*buildNode
}

func (n *buildNode) split(primitiveIndices, indices []uint32, vertices []geometry.Vertex, maxPrimitives, depth int) {
	n.aabb = subMeshAabb(primitiveIndices, indices, vertices)

	if len(primitiveIndices) < maxPrimitives || depth > 16 {
		n.assign(primitiveIndices, indices, vertices)
	} else {
		sp, axis := chooseSplittingPlane(&n.aabb)

		n.axis = axis

		numPids := len(primitiveIndices) / 2
		pids0 := make([]uint32, 0, numPids)
		pids1 := make([]uint32, 0, numPids)

		for _, pi := range primitiveIndices {
			s := triangleSide(vertices[indices[pi + 0]].P, vertices[indices[pi + 1]].P, vertices[indices[pi + 2]].P, sp)
			
			if s == 0 {
				pids0 = append(pids0, pi)
			} else {
				pids1 = append(pids1, pi)
			}
		}

		n.children[0] = new(buildNode)
		n.children[1] = new(buildNode)

		n.children[0].split(pids0, indices, vertices, maxPrimitives, depth + 1)
		n.children[1].split(pids1, indices, vertices, maxPrimitives, depth + 1)
	}
}

func (n *buildNode) assign(primitiveIndices []uint32, indices []uint32, vertices []geometry.Vertex) {
//	n.indices = primitiveIndices

	n.triangles = make([]primitive.Triangle, len(primitiveIndices))

	for i, pi := range primitiveIndices {
		n.triangles[i] = primitive.MakeTriangle(&vertices[indices[pi + 0]], &vertices[indices[pi + 1]], &vertices[indices[pi + 2]])
	}
}


var axis = [...]math.Vector3{ 
	math.MakeVector3(1.0, 0.0, 0.0),
	math.MakeVector3(0.0, 1.0, 0.0),
	math.MakeVector3(0.0, 0.0, 1.0), 
}

func (n *buildNode) intersect(ray *math.OptimizedRay, intersection *primitive.Intersection) bool {
	if !n.aabb.IntersectP(ray) {
		return false
	}

	hit := false

	if n.children[0] != nil {
		c := ray.DirIsNeg[n.axis]

		if n.children[c].intersect(ray, intersection) {
			hit = true
		} 

		if n.children[1 - c].intersect(ray, intersection) {
			hit = true
		}
	} else {
		var ti primitive.Intersection
		var index int

		for i, t := range n.triangles {
			if t.Intersect(ray, &ti.T, &ti.U, &ti.V) {
				if ti.T <= intersection.T {
					*intersection = ti
					index = i
					hit = true
				}
			}
		}

		if hit {
			intersection.Triangle = &n.triangles[index]
			ray.MaxT = intersection.T
		}
	}

	return hit
}

func (n *buildNode) intersectP(ray *math.OptimizedRay) bool {
	if !n.aabb.IntersectP(ray) {
		return false
	}

	if n.children[0] != nil {
		c := ray.DirIsNeg[n.axis]

		if n.children[c].intersectP(ray) {
			return true
		} 

		return n.children[1 - c].intersectP(ray)
	}

	for _, t := range n.triangles {
		if t.IntersectP(ray) {
			return true
		}
	}

	return false
}

func subMeshAabb(primitiveIndices, indices []uint32, vertices []geometry.Vertex) bounding.AABB {
	min := math.MakeVector3( gomath.MaxFloat32,  gomath.MaxFloat32,  gomath.MaxFloat32)
	max := math.MakeVector3(-gomath.MaxFloat32, -gomath.MaxFloat32, -gomath.MaxFloat32)
	
	for _, pi := range primitiveIndices {
		min = triangleMin(vertices[indices[pi + 0]].P, vertices[indices[pi + 1]].P, vertices[indices[pi + 2]].P, min)
		max = triangleMax(vertices[indices[pi + 0]].P, vertices[indices[pi + 1]].P, vertices[indices[pi + 2]].P, max)
	}

	return bounding.MakeAABB(min, max)
}

func triangleMin(a, b, c, x math.Vector3) math.Vector3 {
	return a.Min(b).Min(c).Min(x)
}

func triangleMax(a, b, c, x math.Vector3) math.Vector3 {
	return a.Max(b).Max(c).Max(x)
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

func splittingPlane(aabb *bounding.AABB) (int, float32) {
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

func chooseSplittingPlane(aabb *bounding.AABB) (math.Plane, int) {
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