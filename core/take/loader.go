package take

import (
	pkgfilm "github.com/Opioid/scout/core/rendering/film"
	"github.com/Opioid/scout/core/rendering/integrator"
	"github.com/Opioid/scout/core/rendering/sampler"
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
			take.loadSampler(value)
		case "integrator" == key:
			take.loadIntegrator(value)
		}
	} 

	if (take.Integrator == nil) {
		take.Integrator = integrator.NewWhitted(1)
	}

	if take.Context.Camera == nil {
		return false
	}

	take.Context.Sampler.Resize(math.Vector2i{0, 0}, take.Context.Camera.Film().Dimensions())

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
			if filmNode, ok := value.(map[string]interface{}); ok {
				d := pkgjson.ReadVector2i(filmNode, "dimensions", math.Vector2i{0, 0})
				film = pkgfilm.NewUnfiltered(d)
			}
		}
	}

	var camera pkgcamera.Camera

	switch typestring {
	case "Orthographic":
		camera = pkgcamera.NewOrthographic(dimensions, film)
	case "Perspective":
		camera = pkgcamera.NewPerspective(fov, dimensions, film)
	}

	camera.Transformation().Set(position, math.Vector3{1.0, 1.0, 1.0}, rotation)
	camera.UpdateView()
	take.Context.Camera = camera
}

func (take *Take) loadSampler(s interface{}) {
	samplerNode, ok := s.(map[string]interface{})

	if !ok {
		return
	}

	for key, value := range samplerNode {
		switch key {
		case "Uniform":
			take.Context.Sampler = loadUniformSampler(value)
		}
	}
}

func loadUniformSampler(s interface{}) sampler.Sampler {
	samplerNode, ok := s.(map[string]interface{})

	if !ok {
		return nil
	}

	samplesPerPixel := math.Vector2i{1, 1}

	for key, value := range samplerNode {
		switch key {
		case "samples_per_pixel":
			samplesPerPixel = pkgjson.ParseVector2i(value)
		}
	}

	return sampler.NewUniform(math.Vector2i{}, math.Vector2i{}, samplesPerPixel)
}

func (take *Take) loadIntegrator(i interface{}) {
	integratorNode, ok := i.(map[string]interface{})

	if !ok {
		return
	}

	for key, value := range integratorNode {
		switch key {
		case "Whitted":
			take.Integrator = loadWhittedIntegrator(value)
		case "AO":
			take.Integrator = loadAoIntegrator(value)
		}
	}
}

func loadWhittedIntegrator(i interface{}) integrator.Integrator {
	integratorNode, ok := i.(map[string]interface{})

	if !ok {
		return nil
	}

	bounceDepth := 1

	for key, value := range integratorNode {
		switch key {
		case "bounce_depth":
			bounceDepth = int(value.(float64))
		}
	}

	return integrator.NewWhitted(bounceDepth)
}

func loadAoIntegrator(i interface{}) integrator.Integrator {
	integratorNode, ok := i.(map[string]interface{})

	if !ok {
		return nil
	}

	numSamples := 1
	radius := float32(1.0)

	for key, value := range integratorNode {
		switch key {
		case "num_samples":
			numSamples = int(value.(float64))
		case "radius":
			radius = float32(value.(float64))
		}
	}

	return integrator.NewAo(numSamples, radius)
}