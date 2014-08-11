package take

import (
	pkgfilm "github.com/Opioid/scout/core/rendering/film"
	"github.com/Opioid/scout/core/rendering/sampler"
	"github.com/Opioid/scout/core/scene/camera"
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
		}
	} 

	take.Context.Sampler.Resize(math.Vector2i{0, 0}, take.Context.Camera.Film().Dimensions())

	return true
}

func (take *Take) loadCamera(s interface{}) {
	cameraNode, ok := s.(map[string]interface{})

	if !ok {
		return
	}

	var typestring string
	var position math.Vector3
	var rotation math.Quaternion
	var fov float32
	var dimensions math.Vector2
	var film pkgfilm.Film

	for key, value := range cameraNode {
		switch key {
		case "type":
			typestring = value.(string)
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

	if "Orthographic" == typestring {
		camera := camera.NewOrthographic(dimensions, film)
		camera.Entity.Transformation.Set(position, math.Vector3{1.0, 1.0, 1.0}, rotation)
		take.Context.Camera = camera
	} else if "Perspective" == typestring {
		camera := camera.NewPerspective(fov, dimensions, film)
		camera.Entity.Transformation.Set(position, math.Vector3{1.0, 1.0, 1.0}, rotation)
		take.Context.Camera = camera
	}

	take.Context.Camera.UpdateView()
}

func (take *Take) loadSampler(s interface{}) {
	samplerNode, ok := s.(map[string]interface{})

	if !ok {
		return
	}

	var typestring string
	var samplesPerPixel math.Vector2i

	for key, value := range samplerNode {
		switch key {
		case "type":
			typestring = value.(string)
		case "samples_per_pixel":
			samplesPerPixel = pkgjson.ParseVector2i(value)
		}
	}

	if "Uniform" == typestring {
		take.Context.Sampler = sampler.NewUniform(math.Vector2i{}, math.Vector2i{}, samplesPerPixel)
	}
}