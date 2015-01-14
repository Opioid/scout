package prop

import (
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/scene/material"
	"github.com/Opioid/scout/base/math/bounding"
)

type Prop struct {
	Shape shape.Shape
	AABB bounding.AABB

	Materials []material.Material
}
