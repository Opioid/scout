package scene

import (
	"github.com/Opioid/scout/core/scene/shape"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"io/ioutil"
	"encoding/json"
)

func (scene *Scene) Load(filename string, resourceManager *ResourceManager) error {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	var document interface{}
	if err = json.Unmarshal(data, &document); err != nil {
		return err
	}

	root := document.(map[string]interface{})

	for key, value := range root {
		if "props" == key {
			scene.loadProps(value)
		}
	} 

	return nil
}

func (scene *Scene) loadProps(i interface{}) {
	props, isArray := i.([]interface{})

	if !isArray {
		return 
	}

	for _, prop := range props {
		scene.loadProp(prop)
	}
}

func (scene *Scene) loadProp(i interface{}) {
	propNode, isMap := i.(map[string]interface{})

	if !isMap {
		return
	}

	s, hasShape := propNode["shape"]

	if !hasShape {
		return
	}

	prop := NewProp(loadShape(s))

	for key, value := range propNode {
		switch key {
		case "position":
			prop.Transformation.Position = pkgjson.ParseVector3(value)
		case "scale":
			prop.Transformation.Scale = pkgjson.ParseVector3(value)
		}
	}

	scene.Props = append(scene.Props, prop)
}

func loadShape(i interface{}) shape.Shape {
	shapeNode, isMap := i.(map[string]interface{})

	if !isMap {
		return nil
	}

	t, hasType := shapeNode["type"]

	if !hasType {
		return nil
	}

	switch t {
	case "Sphere":
		return new(shape.Sphere)
	case "Plane":
		return new(shape.Plane)
	default:
		return nil
	}
}