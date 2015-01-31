package take

import (
	pkgfilm "github.com/Opioid/scout/core/rendering/film"
	"github.com/Opioid/scout/core/rendering/film/tonemapping"
	"github.com/Opioid/scout/core/rendering/integrator"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering"
	"github.com/Opioid/scout/core/scene/entity"
	pkgcamera "github.com/Opioid/scout/core/scene/camera"
	"github.com/Opioid/scout/base/math"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

func (take *Take) Load(filename string) bool {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return false
	}

	var document interface{}
	if err = json.Unmarshal(data, &document); err != nil {
		fmt.Println(err)
		return false
	}

	root := document.(map[string]interface{})

	for key, value := range root {
		switch {
		case "scene" == key:
			take.Scene = value.(string)	
		case "camera" == key:
			take.loadCamera(value)
		case "sampler" == key:
			take.Context.Sampler = loadSampler(value)
		case "integrator" == key:
			take.IntegratorFactory = loadIntegratorFactory(value)
		}
	} 

	if (take.IntegratorFactory == nil) {
		take.IntegratorFactory = integrator.NewWhittedFactory(1)
	}

	if take.Context.Camera == nil {
		return false
	}

	if take.Context.Sampler == nil {
		take.Context.Sampler = sampler.NewUniform(math.MakeVector2i(1, 1))
	}

	take.Context.ShutterOpen = 0.0
	take.Context.ShutterClose = 1.0 / 24.0

	return true
}

func (take *Take) loadCamera(c interface{}) {
	cameraNode, ok := c.(map[string]interface{})
	if !ok {
		return
	}

	typestring := "Perspective"
	var typevalue interface{}

	for key, value := range cameraNode {
		typestring = key
		typevalue = value
		break
	}

	settingsNode, ok := typevalue.(map[string]interface{})
	if !ok {
		return
	}

	var position math.Vector3
	rotation := math.MakeIdentityQuaternion()
	lensRadius := float32(0.0)
	focalDistance := float32(0.0)
	shutterSpeed := float32(1.0 / 60.0)	
	fov := float32(60.0)
	var dimensions math.Vector2
	var film pkgfilm.Film
	var animation entity.Animation

	for key, value := range settingsNode {
		switch key {
		case "position":
			position = pkgjson.ParseVector3(value)
		case "rotation":
			rotation = pkgjson.ParseRotationQuaternion(value)
		case "lens_radius":
			lensRadius = float32(value.(float64))
		case "focal_distance":
			focalDistance = float32(value.(float64))
		case "shutter_speed":
			shutterSpeed = float32(value.(float64))			
		case "fov":
			fov = math.DegreesToRadians(float32(value.(float64)))
		case "dimensions":
			dimensions = pkgjson.ParseVector2(value)
		case "film":
			film = loadFilm(value)
		case "keyframes":
			animation = entity.MakeAnimationFromJson(value)
		}
	}

	var camera pkgcamera.Camera

	switch typestring {
	case "Orthographic":
		camera = pkgcamera.NewOrthographic(dimensions, film, shutterSpeed)
	case "Perspective":
		camera = pkgcamera.NewPerspective(lensRadius, focalDistance, fov, dimensions, film, shutterSpeed)
	}

	camera.Entity().SetTransformation(position, math.MakeIdentityVector3(), rotation)
	camera.Entity().Animation = animation	
	camera.UpdateView()
	take.Context.Camera = camera
}

func loadSampler(s interface{}) sampler.Sampler {
	samplerNode, ok := s.(map[string]interface{})
	if !ok {
		return nil
	}

	for key, value := range samplerNode {
		switch key {
		case "Uniform":
			return loadUniformSampler(value)
		case "Stratified":
			return loadStratifiedSampler(value)
		case "Hammersley":
			return loadHammersleySampler(value)				
		case "Scrambled_hammersley":
			return loadScrambledHammersleySampler(value)		
		case "Random":
			return loadRandomSampler(value)			
		}
	}

	return nil
}

func loadUniformSampler(s interface{}) sampler.Sampler {
	samplerNode, ok := s.(map[string]interface{})
	if !ok {
		return nil
	}

	samplesPerPixel := math.MakeVector2i(1, 1)

	for key, value := range samplerNode {
		switch key {
		case "samples_per_pixel":
			samplesPerPixel = pkgjson.ParseVector2i(value)
		}
	}

	return sampler.NewUniform(samplesPerPixel)
}

func loadStratifiedSampler(s interface{}) sampler.Sampler {
	samplerNode, ok := s.(map[string]interface{})
	if !ok {
		return nil
	}

	samplesPerPixel := math.MakeVector2i(1, 1)

	for key, value := range samplerNode {
		switch key {
		case "samples_per_pixel":
			samplesPerPixel = pkgjson.ParseVector2i(value)
		}
	}

	return sampler.NewStratified(samplesPerPixel, nil)
}

func loadHammersleySampler(s interface{}) sampler.Sampler {
	samplerNode, ok := s.(map[string]interface{})
	if !ok {
		return nil
	}

	samplesPerPixel := uint32(1)

	for key, value := range samplerNode {
		switch key {
		case "samples_per_pixel":
			samplesPerPixel = uint32(value.(float64))
		}
	}

	return sampler.NewHammersley(samplesPerPixel)
}

func loadScrambledHammersleySampler(s interface{}) sampler.Sampler {
	samplerNode, ok := s.(map[string]interface{})
	if !ok {
		return nil
	}

	samplesPerPixel := uint32(1)

	for key, value := range samplerNode {
		switch key {
		case "samples_per_pixel":
			samplesPerPixel = uint32(value.(float64))
		}
	}

	return sampler.NewScrambledHammersley(samplesPerPixel, nil)
}

func loadRandomSampler(s interface{}) sampler.Sampler {
	samplerNode, ok := s.(map[string]interface{})
	if !ok {
		return nil
	}

	samplesPerPixel := uint32(1)

	for key, value := range samplerNode {
		switch key {
		case "samples_per_pixel":
			samplesPerPixel = uint32(value.(float64))
		}
	}

	return sampler.NewRandom(samplesPerPixel, nil)
}

func loadFilm(f interface{}) pkgfilm.Film {
	filmNode, ok := f.(map[string]interface{})
	if !ok {
		return nil
	}

	var dimensions math.Vector2i
	var exposure float32
	var tonemapper tonemapping.Tonemapper

	for key, value := range filmNode {
		switch key {
		case "dimensions":
			dimensions = pkgjson.ParseVector2i(value)
		case "exposure":
			exposure = pkgjson.ParseFloat32(value)
		case "tonemapper":
			tonemapper = loadTonemapper(value)
		}
	}

	if (tonemapper == nil) {
		tonemapper = tonemapping.NewIdentity()
	}

	return pkgfilm.NewFiltered(dimensions, exposure, tonemapper)
}

func loadTonemapper(t interface{}) tonemapping.Tonemapper {
	tonemapperNode, ok := t.(map[string]interface{})
	if !ok {
		return nil
	}

	for key, value := range tonemapperNode {
		switch key {
		case "Identity":
			return tonemapping.NewIdentity()
		case "Filmic":
			return loadFilmicTonemapper(value)
		}
	}

	return nil
}

func loadFilmicTonemapper(f interface{}) tonemapping.Tonemapper {
	filmicNode, ok := f.(map[string]interface{})
	if !ok {
		return nil
	}

	w := pkgjson.ReadVector3(filmicNode, "linear_white", math.MakeVector3(1.0, 1.0, 1.0))

	return tonemapping.NewFilmic(w)
}

func loadIntegratorFactory(i interface{}) rendering.IntegratorFactory {
	integratorNode, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	for key, value := range integratorNode {
		switch key {
		case "Whitted":
			return loadWhittedIntegrator(value)
		case "AO":
			return loadAoIntegrator(value)
		case "PT":
			return loadPathtracerIntegrator(value)
		case "PTDL":
			return loadPathtracerDlIntegrator(value)
		case "Normal":
			return loadNormalIntegrator(value)
		}
	}

	return nil
}

func loadWhittedIntegrator(i interface{}) rendering.IntegratorFactory {
	integratorNode, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	maxBounces := uint32(1)

	for key, value := range integratorNode {
		switch key {
		case "max_bounces":
			maxBounces = uint32(value.(float64))
		}
	}

	return integrator.NewWhittedFactory(maxBounces)
}

func loadAoIntegrator(i interface{}) rendering.IntegratorFactory {
	integratorNode, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	numSamples := uint32(8)
	radius := float32(2.0)

	for key, value := range integratorNode {
		switch key {
		case "num_samples":
			numSamples = uint32(value.(float64))
		case "radius":
			radius = float32(value.(float64))
		}
	}

	return integrator.NewAoFactory(numSamples, radius)
}

func loadPathtracerIntegrator(i interface{}) rendering.IntegratorFactory {
	integratorNode, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	maxBounces := uint32(2)

	for key, value := range integratorNode {
		switch key {
		case "max_bounces":
			maxBounces = uint32(value.(float64))
		}
	}

	return integrator.NewPathtracerFactory(maxBounces)
}

func loadPathtracerDlIntegrator(i interface{}) rendering.IntegratorFactory {
	integratorNode, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	maxBounces := uint32(2)

	for key, value := range integratorNode {
		switch key {
		case "max_bounces":
			maxBounces = uint32(value.(float64))
		}
	}

	return integrator.NewPathtracerDlFactory(maxBounces)
}

func loadNormalIntegrator(i interface{}) rendering.IntegratorFactory {
	_/*integratorNode*/, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	return integrator.NewNormalFactory()
}