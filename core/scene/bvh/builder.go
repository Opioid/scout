package bvh

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	_ "fmt"
)

type Builder struct {
//	nodes []miniNode
//	numNodes, currentNode uint32
}


func (b *Builder) Build(props []*prop.StaticProp, maxShapes int, tree *Tree) {
	numFinitShapes := 0

	for _, p := range props {
		if p.Shape.IsFinite() {
			numFinitShapes++
		}
	}

	tree.infiniteProps = make([]*prop.StaticProp, len(props) - numFinitShapes)

	i := 0
	for _, p := range props {
		if !p.Shape.IsFinite() {
			tree.infiniteProps[i] = p
			i++
		}
	}

	indices := make([]uint32, numFinitShapes)

	i = 0
	for pi, p := range props {
		if p.Shape.IsFinite() {
			indices[i] = uint32(pi)
			i++
		}
	}

	root := buildNode{}

	root.split(indices, props, maxShapes)
/*
	b.numNodes = 1
	root.numSubNodes(&b.numNodes)

	fmt.Println(b.numNodes)

	b.nodes = tree.allocateNodes(b.numNodes)

	b.currentNode = 0
	b.serialize(&root)
*/
	tree.root = root
}

func (b *Builder) serialize(node *buildNode) {
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
/*
	n := b.newNode()
	n.aabb = node.aabb
	n.indices = node.indices

	if node.children[0] != nil {
		b.serialize(node.children[0])
		b.serialize(node.children[1])

		n.setHasChildren(true)
	}

	n.setSkipOffset(b.skipOffset())
	*/
}

/*
func (b *Builder) newNode() *miniNode {
	n := &b.nodes[b.currentNode]
	b.currentNode++
	return n
}

func (b *Builder) skipOffset() uint32 {
	if b.currentNode >= b.numNodes {
		return 0
	}

	return b.currentNode
}
*/

type buildNode struct {
	aabb bounding.AABB

	indices []uint32

	children [2]*buildNode
}

func (n *buildNode) split(indices []uint32, props []*prop.StaticProp, maxShapes int) {
	n.aabb = miniaabb(indices, props)

	if len(indices) <= maxShapes {
		n.assign(indices)
	} else {
		n.children[0] = new(buildNode)
		n.children[1] = new(buildNode)

		numProps := len(indices) / 2
		indices0 := make([]uint32, 0, numProps)
		indices1 := make([]uint32, 0, numProps)

		sp := chooseSplittingPlane(&n.aabb)

		for _, i := range indices {
			mib := sp.Behind(props[i].AABB.Min)
			mab := sp.Behind(props[i].AABB.Max)
			if mib && mab {
				indices0 = append(indices0, uint32(i))
			} else if !mib && !mab {
				indices1 = append(indices1, uint32(i))
			} else {
				indices0 = append(indices0, uint32(i))
				indices1 = append(indices1, uint32(i))
			}
		}

		n.children[0].split(indices0, props, maxShapes)
		n.children[1].split(indices1, props, maxShapes)
	}
}

func (n *buildNode) assign(indices []uint32) {
	n.indices = indices
}

func (n *buildNode) numSubNodes(num *uint32) {
	if n.children[0] != nil {
		*num += 2

		n.children[0].numSubNodes(num);
		n.children[1].numSubNodes(num);
	}
}

func (n *buildNode) intersect(ray *math.OptimizedRay, props []*prop.StaticProp, intersection *prop.Intersection) bool {
	if !n.aabb.Intersect(ray) {
		return false
	}

	hit := false

	if n.children[0] != nil {
		sd0 := n.children[0].aabb.Position().SquaredDistance(ray.Origin)
		sd1 := n.children[1].aabb.Position().SquaredDistance(ray.Origin)

		var c0, c1 int

		if sd0 <= sd1 {
			c0 = 0
			c1 = 1
		} else {
			c0 = 1
			c1 = 0
		}

		if n.children[c0].intersect(ray, props, intersection) {
			hit = true
		} 

		if n.children[c1].intersect(ray, props, intersection) {
			hit = true
		}
	} else {
		for _, i := range n.indices {
			p := props[i]
			if p.Intersect(ray, intersection) {
				intersection.Prop = &p.Prop
				hit = true
			}
		}
	}

	return hit
}

func (n *buildNode) intersectP(ray *math.OptimizedRay, props []*prop.StaticProp) bool {
	if !n.aabb.Intersect(ray) {
		return false
	}

	if n.children[0] != nil {
		if n.children[0].intersectP(ray, props) {
			return true
		} 

		return n.children[1].intersectP(ray, props)
	}

	for _, i := range n.indices {
		if props[i].IntersectP(ray) {
			return true
		}
	}

	return false
}

func miniaabb(indices []uint32, props []*prop.StaticProp) bounding.AABB {
	b := bounding.MakeAABB()

	for _, i := range indices {
		b = b.Merge(&props[i].AABB)
	}

	return b
}

func chooseSplittingPlane(aabb *bounding.AABB) math.Plane {
	position := aabb.Position()
	halfsize := aabb.Halfsize()

	if halfsize.X >= halfsize.Y && halfsize.X >= halfsize.Z {
		return math.MakePlane(math.Vector3{1.0, 0.0, 0.0}, position)
	} else if halfsize.Y >= halfsize.X && halfsize.Y >= halfsize.Z {
		return math.MakePlane(math.Vector3{0.0, 1.0, 0.0}, position)
	} else {
		return math.MakePlane(math.Vector3{0.0, 0.0, 1.0}, position)
	}
}