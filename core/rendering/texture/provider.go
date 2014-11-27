package texture

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/rendering/color"
	_ "github.com/Opioid/scout/base/file"
	"os"
	goimage "image"
	_ "image/jpeg"
	"runtime"
	"sync"
	"fmt"
)

type Provider struct {

}

func (p *Provider) Load2D(filename string, treatAsLinear bool) *Texture2D {
	fi, err := os.Open(filename)

	if err != nil {
		fmt.Printf("%s could not be loaded", filename)
		return nil
	}

	defer fi.Close()

	// fmt.Println(file.QueryFileType(fi))

	sourceImage, _, err := goimage.Decode(fi)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	dimensions := math.MakeVector2i(int32(sourceImage.Bounds().Max.X), int32(sourceImage.Bounds().Max.Y))

	texture := NewTexture2D(dimensions, 1)

	numTaks := int32(runtime.GOMAXPROCS(0))

	a := dimensions.Y / numTaks

	start := math.MakeVector2i(0, 0)
	end   := math.MakeVector2i(dimensions.X, a)

	wg := sync.WaitGroup{}

	for i := int32(0); i < numTaks; i++ {
		wg.Add(1)

		go func (start, end math.Vector2i) {
			if treatAsLinear {
				processLinear(start, end, sourceImage, &texture.Image.Buffers[0])
			} else {
				processSrgbToLinear(start, end, sourceImage, &texture.Image.Buffers[0])
			}

			wg.Done()
		}(start, end)

		start.Y += a

		if i == numTaks - 2 {
			end.Y = dimensions.Y
		} else {
			end.Y += a
		}
	}

	wg.Wait()

	return texture
}

func processSrgbToLinear(start, end math.Vector2i, source goimage.Image, target *Buffer) {
	max := float32(0xFFFF)

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			r, g, b, a := source.At(int(x), int(y)).RGBA()

			target.Set(x, y, math.MakeVector4(
				color.SrgbToLinear(float32(r) / max), 
				color.SrgbToLinear(float32(g) / max), 
				color.SrgbToLinear(float32(b) / max), 
				float32(a) / max,
			))
		}
	}
}

func processLinear(start, end math.Vector2i, source goimage.Image, target *Buffer) {
	max := float32(0xFFFF)

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			r, g, b, a := source.At(int(x), int(y)).RGBA()

			target.Set(x, y, math.MakeVector4(
				float32(r) / max, 
				float32(g) / max, 
				float32(b) / max, 
				float32(a) / max,
			))
		}
	}
}