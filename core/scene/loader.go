package scene

import (
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"io/ioutil"
	"encoding/json"
	_ "fmt"
)

type Loader struct {
	scene *Scene
	resourceManager *ResourceManager
}

func NewLoader(scene *Scene, resourceManager *ResourceManager) *Loader {
	loader := new(Loader)

	loader.scene = scene
	loader.resourceManager = resourceManager

	return loader
}

func (loader *Loader) Load(filename string) error {
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
		if "entities" == key {
			loader.loadEntities(value)
		} else if "static_props" == key {
			loader.loadStaticProps(value)
		} 
	} 

	loader.scene.Compile()

	return nil
}

func (loader *Loader) loadEntities(i interface{}) {
	entities, ok := i.([]interface{})

	if !ok {
		return 
	}

	for _, entity := range entities {
		entityNode, ok := entity.(map[string]interface{})

		if !ok {
			continue
		}

		classNode, ok := entityNode["class"]

		if !ok {
			continue
		}

		className := classNode.(string)

		if "Light" == className {
			loader.loadLight(entityNode)
		}
	}
}

func (loader *Loader) loadLight(i interface{}) {
	lightNode, ok := i.(map[string]interface{})

	if !ok {
		return
	}

	typeNode, ok := lightNode["type"]

	if !ok {
		return
	}

	typeName := typeNode.(string)

	var l *light.Light

	if "Directional" == typeName {
		l = loader.scene.CreateLight(light.Directional)
	}

	var position, scale math.Vector3
	var rotation math.Quaternion

	for key, value := range lightNode {
		switch key {
		case "color":
			l.Color = pkgjson.ParseVector3(value)

		case "position":
			position = pkgjson.ParseVector3(value)
		case "scale":
			scale = pkgjson.ParseVector3(value)
		case "rotation":
			rotation = pkgjson.ParseRotationQuaternion(value)
		}
	}

	l.Transformation.Set(position, scale, rotation)
/*
	s, ok := propNode["shape"]

	if !ok {
		return
	}

	shape := loader.resourceManager.LoadShape(s.(string))

	if shape == nil {
		return
	}

	m, ok := propNode["material"]

	if !ok {
		return
	}

	material := loader.resourceManager.LoadMaterial(m.(string))

	if material == nil {
		return
	}

	var position math.Vector3
	var scale math.Vector3
	rotation := math.MakeIdentityMatrix3x3()

	for key, value := range propNode {
		switch key {
		case "position":
			position = pkgjson.ParseVector3(value)
		case "scale":
			scale = pkgjson.ParseVector3(value)
		case "rotation":
			rotation = pkgjson.ParseRotationMatrix(value)
		}
	}

	prop := NewStaticProp(shape, material)

	prop.SetWorldTransformation(position, scale, &rotation)

	loader.scene.StaticProps = append(loader.scene.StaticProps, prop)
	*/
}

func (loader *Loader) loadStaticProps(i interface{}) {
	props, ok := i.([]interface{})

	if !ok {
		return 
	}

	for _, prop := range props {
		loader.loadStaticProp(prop)
	}
}

func (loader *Loader) loadStaticProp(i interface{}) {
	propNode, ok := i.(map[string]interface{})

	if !ok {
		return
	}

	s, ok := propNode["shape"]

	if !ok {
		return
	}

	shape := loader.resourceManager.LoadShape(s.(string))

	if shape == nil {
		return
	}

	m, ok := propNode["material"]

	if !ok {
		return
	}

	material := loader.resourceManager.LoadMaterial(m.(string))

	if material == nil {
		return
	}

	var position, scale math.Vector3
	rotation := math.MakeIdentityQuaterion()

	for key, value := range propNode {
		switch key {
		case "position":
			position = pkgjson.ParseVector3(value)
		case "scale":
			scale = pkgjson.ParseVector3(value)
		case "rotation":
			rotation = pkgjson.ParseRotationQuaternion(value)
		}
	}

	prop := NewStaticProp(shape, material)

	prop.Transformation.Set(position, scale, rotation)

	loader.scene.StaticProps = append(loader.scene.StaticProps, prop)
}