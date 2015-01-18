package bvh

import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/shape/triangle/primitive"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	_ "math"
	"fmt"
)

type Builder1 struct {
	nodes []node
	numNodes, currentNode uint32
}

func (b *Builder1) Build(triangles []primitive.IndexTriangle, vertices []geometry.Vertex, maxPrimitives int) Tree1 {
	primitiveIndices := make([]uint32, len(triangles))

	for i := range primitiveIndices {
		primitiveIndices[i] = uint32(i)
	}

	outTriangles := make([]primitive.Triangle, 0, len(triangles))

	root := buildNode1{}
	root.split(primitiveIndices, triangles, vertices, maxPrimitives, 0, &outTriangles)

//	return Tree1{root: root, Triangles: outTriangles}

	tree := Tree1{}

	tree.Triangles = outTriangles

	b.numNodes = 1
	root.numSubNodes(&b.numNodes)

	fmt.Println(b.numNodes)

	b.nodes = tree.allocateNodes(b.numNodes)

	b.currentNode = 0
	b.serialize(&root)

	return tree
}

func (b *Builder1) serialize(node *buildNode1) {
/*
	AABB_node* n = get_new_node();

	n->aabb_ = node->aabb;
	n->props_.swap(node->props);

	if (node->child0)
	{
		serialize(node->child0);

		serialize(node->child1);

		n->set_has_children(true);
	}

	n->set_skip_node(skip_node());
*/

	n := b.newNode()
	n.aabb = node.aabb
	n.startIndex = node.startIndex
	n.endIndex = node.endIndex
	n.axis = node.axis

	if node.children[0] != nil {
		b.serialize(node.children[0])

		n.setRightChild(b.currentNodeIndex())

		b.serialize(node.children[1])

		n.setHasChildren(true)
	}
}

func (b *Builder1) newNode() *node {
	n := &b.nodes[b.currentNode]
	b.currentNode++
	return n
}

func (b *Builder1) currentNodeIndex() uint32 {
	return b.currentNode
}


type buildNode1 struct {
	aabb bounding.AABB

	startIndex, endIndex uint32

	children [2]*buildNode1

	axis int8
}

func (n *buildNode1) split(primitiveIndices []uint32, triangles []primitive.IndexTriangle, vertices []geometry.Vertex, maxPrimitives, depth int, 
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

		primitiveIndices = nil

		n.children[0] = new(buildNode1)
		n.children[1] = new(buildNode1)

		n.children[0].split(pids0, triangles, vertices, maxPrimitives, depth + 1, outTriangles)

		pids0 = nil

		n.children[1].split(pids1, triangles, vertices, maxPrimitives, depth + 1, outTriangles)
	}
}

func (n *buildNode1) assign(primitiveIndices []uint32, triangles []primitive.IndexTriangle, vertices []geometry.Vertex, 
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

func (n *buildNode1) numSubNodes(num *uint32) {
	if n.children[0] != nil {
		*num += 2

		n.children[0].numSubNodes(num);
		n.children[1].numSubNodes(num);
	}
}

func (n *buildNode1) intersect(ray *math.OptimizedRay, triangles []primitive.Triangle, intersection *primitive.Intersection) bool {
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

func (n *buildNode1) intersectP(ray *math.OptimizedRay, triangles []primitive.Triangle) bool {
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