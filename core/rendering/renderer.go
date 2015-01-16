package rendering

import (
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/progress"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/random"
	"sync"
	_ "fmt"
)

type Renderer struct {
	IntegratorFactory IntegratorFactory

	tileSize math.Vector2i
	currentPixel math.Vector2i
}

func (r *Renderer) Render(scene *pkgscene.Scene, context *Context, numThreads uint32, progressor progress.Sink) {
	dimensions := context.Camera.Film().Dimensions()

	shutterClose := context.ShutterOpen + context.Camera.ShutterSpeed()

	r.currentPixel = math.MakeVector2i(0, 0)
	r.tileSize     = math.MakeVector2i(32, 32)

	numTiles := int(float32(dimensions.X) / float32(r.tileSize.X) + 0.5) * int(float32(dimensions.Y) / float32(r.tileSize.Y) + 0.5)
	tiles := make(chan tile, numTiles)

	for {
		tiles <- tile{r.currentPixel, r.currentPixel.Add(r.tileSize).Min(dimensions)}

		if !r.advanceCurrentPixel(dimensions) {
			break
		}		
	}

	close(tiles)

	progressor.Start(numTiles)

	wg := sync.WaitGroup{}

	for i := uint32(0); i < numThreads; i++ {

		wg.Add(1)

		go func () {	
			rng := random.Generator{}
			rng.Seed(i + 0, i + 1, i + 2, i + 3)

			worker := makeWorker(r.IntegratorFactory.New(&rng))

			sampler := context.Sampler.Clone(&rng)

			for tile := range tiles {
				worker.render(scene, context.Camera, context.ShutterOpen, shutterClose, tile.start, tile.end, sampler)

				progressor.Tick()
			}

			wg.Done()
		}()
	}

	wg.Wait()
	progressor.End()
}

func (r *Renderer) advanceCurrentPixel(dimensions math.Vector2i) bool {
	r.currentPixel.X += r.tileSize.X

	if r.currentPixel.X >= dimensions.X {
		r.currentPixel.X = 0
		r.currentPixel.Y += r.tileSize.Y
	}

	if r.currentPixel.Y >= dimensions.Y {
		return false
	}

	return true
}

type tile struct {
	start, end math.Vector2i
}