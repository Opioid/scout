package scene

import (
	"github.com/Opioid/scout/core/scene/shape"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

func (scene *Scene) Load(filename string) bool {
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
		if "shapes" == key {
			scene.loadShapes(value)
		}
	} 

	return true
}

func (scene *Scene) loadShapes(s interface{}) {
	shapes, isArray := s.([]interface{})

	if !isArray {
		return 
	}

	for _, shape := range shapes {
		scene.loadShape(shape)
	}
}

func (scene *Scene) loadShape(s interface{}) {
	shapeNode, isMap := s.(map[string]interface{})

	if !isMap {
		return
	}

	t, hasType := shapeNode["type"]

	if !hasType {
		return
	}

	switch t {
	case "Sphere":
		sphere := loadSphere(shapeNode)
		scene.Shapes = append(scene.Shapes, sphere)
	case "Plane":
		scene.Shapes = append(scene.Shapes, new(shape.Plane))
	}
}

func loadSphere(s map[string]interface{}) *shape.Sphere {
	sphere := new(shape.Sphere)

	for key, value := range s {
		switch key {
		case "position":
			sphere.Position = pkgjson.ParseVector3(value)
		case "radius":
			sphere.Radius = float32(value.(float64))
		}
	}

	return sphere
}