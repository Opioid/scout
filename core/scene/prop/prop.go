package prop

import (
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/scene/material"
	"github.com/Opioid/scout/base/math/bounding"
)

type Prop struct {
	Shape shape.Shape

	// This is intended to hold the world space AABB.
	// Therefore it actually only makes sense for static props because it is dependent on time.
	// Maybe move there.
	AABB bounding.AABB	

	Materials []material.Material
}
