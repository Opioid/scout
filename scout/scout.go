/*
package main

import (
	"github.com/Opioid/math32"
	_ "github.com/Opioid/scout/base/math"
	_ "math"
	"time"
	"fmt"
)

func main() {
	fmt.Println("It runs...")

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
	 "github.com/davecheney/profile"
)

func main() {
	cfg := profile.Config {
		CPUProfile:     true,
		MemProfile:  	false,
		ProfilePath:    ".",  // store profiles in current directory
		NoShutdownHook: true, // do not hook SIGINT
	}

	// p.Stop() must be called before the program exits to
	defer profile.Start(&cfg).Stop()

	numWorkers := /*uint32(1)//*/uint32(runtime.NumCPU()) - 0

	fmt.Printf("#Threads %d\n", numWorkers)
	runtime.GOMAXPROCS(int(numWorkers))


	take := take.Take{}

	takename := "../data/takes/cornell.take"

	if !take.Load(takename) {
		fmt.Printf("Take \"%v\" could not be loaded.\n", takename)
		return
	}

	resourceManager := resource.NewManager(numWorkers)

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
	renderer.Render(&scene, &take.Context, numWorkers, progressor)

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