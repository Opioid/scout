package main

import (
	"github.com/Opioid/scout/core/rendering"
	pkgscene "github.com/Opioid/scout/core/scene"
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

	if !take.Load("../data/takes/test.take") {
		fmt.Println("Take could not be loaded")
	}

	resourceManager := pkgscene.NewResourceManager()

	scene := pkgscene.Scene{}

	fmt.Printf("Loading...")
	loadStart := time.Now()

	sceneLoader := pkgscene.NewLoader(&scene, resourceManager)

	if err := sceneLoader.Load(take.Scene); err != nil {
		fmt.Printf("Scene could not be loaded: %s\n",err)
	}

	loadDuration := time.Since(loadStart)
	seconds := float64(loadDuration.Nanoseconds()) / 1000000000.0
	fmt.Printf("     (%fs)\n", seconds)

	renderer := rendering.Renderer{BounceDepth: 1}

	fmt.Printf("Rendering...")
	renderStart := time.Now()

	renderer.Render(&scene, &take.Context)

	renderDuration := time.Since(renderStart)
	seconds = float64(renderDuration.Nanoseconds()) / 1000000000.0
	fmt.Printf("   (%fs)\n", seconds)

	fmt.Printf("Saving...")
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
	fmt.Printf("      (%fs)\n", seconds)
}