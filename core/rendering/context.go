package rendering

import (
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/scene/camera"
)

type Context struct {
	Camera camera.Camera
	Sampler sampler.Sampler
	ShutterOpen, ShutterClose float32
}