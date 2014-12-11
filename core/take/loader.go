package take

import (
	pkgfilm "github.com/Opioid/scout/core/rendering/film"
	"github.com/Opioid/scout/core/rendering/film/tonemapping"
	"github.com/Opioid/scout/core/rendering/integrator"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/rendering"
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
		take.IntegratorFactory = integrator.NewWhittedFactory(1, 16)
	}

	if take.Context.Camera == nil {
		return false
	}

	take.Context.Sampler.Resize(math.MakeVector2i(0, 0), take.Context.Camera.Film().Dimensions())

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
	var rotation math.Quaternion
	var fov float32
	var dimensions math.Vector2
	var film pkgfilm.Film

	for key, value := range settingsNode {
		switch key {
		case "position":
			position = pkgjson.ParseVector3(value)
		case "rotation":
			rotation = pkgjson.ParseRotationQuaternion(value)
		case "fov":
			fov = math.DegreesToRadians(float32(value.(float64)))
		case "dimensions":
			dimensions = pkgjson.ParseVector2(value)
		case "film":
				//w := pkgjson.ReadVector3(filmNode, "linear_white", math.Vector3{1.0, 1.0, 1.0})
				film = loadFilm(value)
		//	}
		}
	}

	var camera pkgcamera.Camera

	switch typestring {
	case "Orthographic":
		camera = pkgcamera.NewOrthographic(dimensions, film)
	case "Perspective":
		camera = pkgcamera.NewPerspective(fov, dimensions, film)
	}

	camera.Transformation().Set(position, math.MakeVector3(1.0, 1.0, 1.0), rotation)
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
		case "Quincunx":
			return loadQuincunxSampler(value)
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

	return sampler.NewUniform(math.Vector2i{}, math.Vector2i{}, samplesPerPixel)
}

func loadQuincunxSampler(s interface{}) sampler.Sampler {
	return sampler.NewQuincunx(math.Vector2i{}, math.Vector2i{})
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
		}
	}

	return nil
}

func loadWhittedIntegrator(i interface{}) rendering.IntegratorFactory {
	integratorNode, ok := i.(map[string]interface{})

	if !ok {
		return nil
	}

	bounceDepth := uint32(1)
	maxLightSamples := uint32(16)

	for key, value := range integratorNode {
		switch key {
		case "bounce_depth":
			bounceDepth = uint32(value.(float64))
		case "max_light_samples":
			maxLightSamples = uint32(value.(float64))
		}
	}

	return integrator.NewWhittedFactory(bounceDepth, maxLightSamples)
}

func loadAoIntegrator(i interface{}) rendering.IntegratorFactory {
	integratorNode, ok := i.(map[string]interface{})

	if !ok {
		return nil
	}

	numSamples := uint32(1)
	radius := float32(1.0)

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