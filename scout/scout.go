package main

import (
	"github.com/Opioid/scout/core/rendering"
	"github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/take"
	"fmt"
	"os"
	"image/png"
)

func main() {

	take := take.Take{}

	if !take.Load("../data/takes/test.take") {
		fmt.Println("Take could not be loaded")
	}

	fmt.Println(take)

	resourceManager := scene.NewResourceManager()

	scene := scene.Scene{}

	if err := scene.Load(take.Scene, resourceManager); err != nil {
		fmt.Printf("Scene could not be loaded: %s\n",err)
	}

	dimensions := take.Camera.Film().Dimensions
	buffer := rendering.NewPixelBuffer(dimensions)

	context := &rendering.Context{}
	context.Camera = take.Camera
	context.Target = buffer;

	renderer := rendering.Renderer{}

	renderer.Render(&scene, context)

	fo, err := os.Create("output.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()
	
	image := buffer.RGBA()

	png.Encode(fo, image)
}