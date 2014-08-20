package scene

import (
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/base/math"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"io/ioutil"
	"encoding/json"
	_ "fmt"
)

type Loader struct {
	scene *Scene
	resourceManager *ResourceManager

	disk   *shape.Disk
	plane  *shape.Plane
	sphere *shape.Sphere
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

		switch className {
		case "Light":
			loader.loadLight(entityNode)
		case "Complex":
			loader.loadComplex(entityNode)
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

	typename := typeNode.(string)

	var l light.Light

	switch typename {
	case "Directional":
		l = light.NewDirectional()
	case "Point":
		l = light.NewPoint()
	}

	var position math.Vector3
	scale := math.MakeIdentityVector3()
	rotation := math.MakeIdentityQuaterion()

	for key, value := range lightNode {
		switch key {
		case "color":
			l.SetColor(pkgjson.ParseVector3(value))
		case "lumen":
			l.SetLumen(pkgjson.ParseFloat32(value))
		case "position":
			position = pkgjson.ParseVector3(value)
		case "scale":
			scale = pkgjson.ParseVector3(value)
		case "rotation":
			rotation = pkgjson.ParseRotationQuaternion(value)
		}
	}

	l.Entity().Transformation.Set(position, scale, rotation)

	loader.scene.AddLight(l)
}

func (loader *Loader) loadComplex(i interface{}) {
	complexNode, ok := i.(map[string]interface{})

	if !ok {
		return
	}

	typeNode, ok := complexNode["type"]

	if !ok {
		return
	}

	typename := typeNode.(string)

	c := loader.scene.CreateComplex(typename)

	c.Init(loader.scene)
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

	shape := loader.loadShape(s)

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
	scale := math.MakeIdentityVector3()
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

	prop.SetTransformation(position, scale, rotation)

	loader.scene.StaticProps = append(loader.scene.StaticProps, prop)
}

func (loader *Loader) loadShape(i interface{}) shape.Shape {
	shapeNode, ok := i.(map[string]interface{})

	if !ok {
		return nil
	}

	if t, ok := shapeNode["type"]; ok {
		typename := t.(string)
		switch typename {
		case "Disk":
			if loader.disk == nil {
				loader.disk = shape.NewDisk()
			} 
			return loader.disk
		case "Plane":
			if loader.plane == nil {
				loader.plane = shape.NewPlane()
			} 
			return loader.plane
		case "Sphere":
			if loader.sphere == nil {
				loader.sphere = shape.NewSphere()
			}
			return loader.sphere
		}
	} else if f, ok := shapeNode["file"]; ok {
		file := f.(string)
		return loader.resourceManager.LoadShape(file)
	}

	return nil
}