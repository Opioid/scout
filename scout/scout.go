/*
package main

import (
	_ "github.com/Opioid/math32"
	"github.com/Opioid/scout/base/math"
	_ "math"
	_ "time"
	"fmt"
)

func main() {
	fmt.Println("It runs...")
	
	a := math.MakeVector3(0.0, 1.0, 2.0)
	b := math.MakeVector3(3.0, 4.0, 5.0)

	t := float32(0.25)

	{
		r := math.Vector3Lerp(&a, &b, t)

		fmt.Println(r)
	}

	{
		r := a.Lerp(b, t)

		fmt.Println(r)
	}
}
*/

package main

import (
	"github.com/Opioid/rgbe"
	"github.com/Opioid/scout/scout/complex"
	"github.com/Opioid/scout/core/rendering"
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/resource"
	"github.com/Opioid/scout/core/take"
	"github.com/Opioid/scout/core/progress"
	"runtime"
	"os"
	"time"
	"fmt"
	"image/png"
//	 "github.com/davecheney/profile"
)

func main() {
//	cfg := profile.Config {
//		CPUProfile:     true,
//		ProfilePath:    ".",  // store profiles in current directory
//		NoShutdownHook: true, // do not hook SIGINT
//	}

	// p.Stop() must be called before the program exits to
//	defer profile.Start(&cfg).Stop()

	fmt.Printf("#Cores %d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	take := take.Take{}

	takename := "../data/takes/motion_test.take"

	if !take.Load(takename) {
		fmt.Printf("Take \"%v\" could not be loaded.\n", takename)
		return
	}

	resourceManager := resource.NewManager()

	scene := pkgscene.Scene{}
	scene.Init()

	complex.Init(&scene)

	fmt.Printf("Loading...\n")
	loadStart := time.Now()

	sceneLoader := pkgscene.NewLoader(&scene, resourceManager)

	if err := sceneLoader.Load(take.Scene); err != nil {
		fmt.Printf("Scene could not be loaded: %s\n", err)
		return
	}

	loadDuration := time.Since(loadStart)
	seconds := float64(loadDuration.Nanoseconds()) / 1000000000.0
	fmt.Printf("(%fs)\n", seconds)

	renderer := rendering.Renderer{}
	renderer.IntegratorFactory = take.IntegratorFactory

	fmt.Printf("Rendering...\n")
	renderStart := time.Now()

	progressor := progress.NewStdout()
	renderer.Render(&scene, &take.Context, progressor)

	renderDuration := time.Since(renderStart)
	seconds = float64(renderDuration.Nanoseconds()) / 1000000000.0
	fmt.Printf("(%fs)\n", seconds)

	fmt.Printf("Saving...\n")
	saveStart := time.Now()

	film := take.Context.Camera.Film()

	{
		fo, err := os.Create("output.png")

		defer fo.Close()

		if err != nil {
			panic(err)
		}

		image := film.RGBA()

		png.Encode(fo, image)
	}

	{
		fo, err := os.Create("output.hdr")

		defer fo.Close()

		if err != nil {
			panic(err)
		}

		data := film.Float32x3()

		dimensions := film.Dimensions()

		if err := rgbe.Encode(fo, int(dimensions.X), int(dimensions.Y), data); err != nil {
			panic(err)
		}
	}

	saveDuration := time.Since(saveStart)
	seconds = float64(saveDuration.Nanoseconds()) / 1000000000.0
	fmt.Printf("(%fs)\n", seconds)
}