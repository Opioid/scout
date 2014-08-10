package texture

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/rendering/color"
	"os"
	goimage "image"
	_ "image/jpeg"
	"fmt"
)

type Provider struct {

}

func (p *Provider) Load2D(filename string) *Texture2D {
	fi, err := os.Open(filename)

	if err != nil {
		return nil
	}

	defer fi.Close()

	image, _, err := goimage.Decode(fi)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	texture := NewTexture2D(1)

	dimensions := image.Bounds().Max
	texture.images[0].resize(math.Vector2i{dimensions.X, dimensions.Y})

	max := float32(0xFFFF)

	for y := 0; y < dimensions.Y; y++ {
		for x := 0; x < dimensions.X; x++ {
			r, g, b, a := image.At(x, y).RGBA()

			texture.images[0].set(x, y, math.Vector4{color.SrgbToLinear(float32(r) / max), 
													 color.SrgbToLinear(float32(g) / max), 
													 color.SrgbToLinear(float32(b) / max), 
													 float32(a) / max})
		}
	}

	return texture
}