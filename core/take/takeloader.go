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
	cameraNode, isMap := s.(map[string]interface{})

	if !isMap {
		return
	}

	var position math.Vector3
	var dimensions math.Vector2i
	var film camera.Film

	for key, value := range cameraNode {
		switch key {
		case "position":
			position = pkgjson.ParseVector3(value)
		case "dimensions":
			dimensions = pkgjson.ParseVector2i(value)
		case "film":
			if filmNode, ok := value.(map[string]interface{}); ok {
				film.Dimensions = pkgjson.ReadVector2i(filmNode, "dimensions", math.Vector2i{0, 0})
			}
		}
	}
	
	camera := camera.NewOrthographic(dimensions, film)
	camera.Entity.Position = position

	take.Camera = camera
}