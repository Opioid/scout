package main

import (
	"github.com/Opioid/scout/core/rendering"
	pkgscene "github.com/Opioid/scout/core/scene"
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

	resourceManager := pkgscene.NewResourceManager()

	scene := pkgscene.Scene{}

	sceneLoader := pkgscene.NewLoader(&scene, resourceManager)

	if err := sceneLoader.Load(take.Scene); err != nil {
		fmt.Printf("Scene could not be loaded: %s\n",err)
	}

	dimensions := take.Camera.Film().Dimensions
	buffer := rendering.NewPixelBuffer(dimensions)

	context := &rendering.Context{}
	context.Camera = take.Camera
	context.Target = buffer;

	renderer := rendering.Renderer{BounceDepth: 1}

	renderer.Render(&scene, context)

	fo, err := os.Create("output.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()
	
	image := buffer.RGBA()

	png.Encode(fo, image)
}