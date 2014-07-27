package take

import (
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
		}
	} 

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
	var film camera.Film

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
				film.Dimensions = pkgjson.ReadVector2i(filmNode, "dimensions", math.Vector2i{0, 0})
			}
		}
	}

	if "Orthographic" == typestring {
		camera := camera.NewOrthographic(dimensions, film)
		camera.Entity.Transformation.Position = position
		camera.Entity.Transformation.Rotation = rotation

		take.Camera = camera
	} else if "Perspective" == typestring {
		camera := camera.NewPerspective(fov, dimensions, film)
		camera.Entity.Transformation.Position = position
		camera.Entity.Transformation.Rotation = rotation

		take.Camera = camera
	}

	take.Camera.UpdateView()
}