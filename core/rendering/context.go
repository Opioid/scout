package rendering

import (
	"github.com/Opioid/scout/core/scene/camera"
)

type Context struct {
	Camera camera.Camera
	Target *PixelBuffer
}