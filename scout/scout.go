package main

import (
	"github.com/Opioid/scout/scout/complex"
	"github.com/Opioid/scout/core/rendering"
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/resource"
	"github.com/Opioid/scout/core/take"
	"runtime"
	"os"
	"time"
	"fmt"
	"image/png"
)

func main() {
	fmt.Printf("#Cores %d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	take := take.Take{}

	if !take.Load("../data/takes/bvh_test.take") {
		fmt.Println("Take could not be loaded")
	}

	resourceManager := resource.NewManager()

	scene := pkgscene.Scene{}
	scene.Init()

	complex.Init(&scene)

	fmt.Printf("Loading...\n")
	loadStart := time.Now()

	sceneLoader := pkgscene.NewLoader(&scene, resourceManager)

	if err := sceneLoader.Load(take.Scene); err != nil {
		fmt.Printf("Scene could not be loaded: %s\n",err)
	}

	loadDuration := time.Since(loadStart)
	seconds := float64(loadDuration.Nanoseconds()) / 1000000000.0
	fmt.Printf("(%fs)\n", seconds)

	renderer := rendering.Renderer{BounceDepth: 1}

	fmt.Printf("Rendering...\n")
	renderStart := time.Now()

	renderer.Render(&scene, &take.Context)

	renderDuration := time.Since(renderStart)
	seconds = float64(renderDuration.Nanoseconds()) / 1000000000.0
	fmt.Printf("(%fs)\n", seconds)

	fmt.Printf("Saving...\n")
	saveStart := time.Now()

	image := take.Context.Camera.Film().RGBA()

	fo, err := os.Create("output.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	png.Encode(fo, image)

	saveDuration := time.Since(saveStart)
	seconds = float64(saveDuration.Nanoseconds()) / 1000000000.0
	fmt.Printf("(%fs)\n", seconds)
}