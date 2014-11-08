package scene

import (
	"github.com/Opioid/scout/core/scene/surrounding"
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/resource"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"io/ioutil"
	"encoding/json"
	_ "fmt"
)

type Loader struct {
	scene *Scene
	resourceManager *resource.Manager

	disk   *shape.Disk
	plane  *shape.Plane
	sphere *shape.Sphere
}

func NewLoader(scene *Scene, resourceManager *resource.Manager) *Loader {
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
		switch key {
		case "surrounding":
			loader.loadSurrounding(value)
		case "entities":
			loader.loadEntities(value)
		case "static_props":
			loader.loadStaticProps(value)
		} 
	} 

	if loader.scene.Surrounding == nil {
		loader.scene.Surrounding = surrounding.NewUniform(math.MakeVector3(0.0, 0.0, 0.0))
	}

	loader.scene.Compile()

	return nil
}

func (loader *Loader) loadSurrounding(i interface{}) {
	surroundingNode, ok := i.(map[string]interface{})

	if !ok {
		return
	}

	typeNode, ok := surroundingNode["type"]

	if !ok {
		return
	}

	typename := typeNode.(string)

	switch typename {
	case "Uniform": 
		color := pkgjson.ReadVector3(surroundingNode, "color", math.MakeVector3(0.0, 0.0, 0.0))
		loader.scene.Surrounding = surrounding.NewUniform(color)

	case "Textured":

		t, ok := surroundingNode["texture"]

		if !ok {
			return
		}

		textureNode, ok := t.(map[string]interface{})

		filename := textureNode["file"].(string)

		if sphericalTexture := loader.resourceManager.LoadTexture2D(filename); sphericalTexture != nil {
			sampler := texture.NewSamplerSpherical_bilinear(sphericalTexture)
			loader.scene.Surrounding = surrounding.NewSphere(sampler)
		}
	}
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


	color := math.MakeVector3(1.0, 1.0, 1.0)
	lumen := float32(1.0)
	var position math.Vector3
	scale := math.MakeIdentityVector3()
	rotation := math.MakeIdentityQuaternion()
	radius := float32(1.0)

	for key, value := range lightNode {
		switch key {
		case "color":
			color = pkgjson.ParseVector3(value)
		case "lumen":
			lumen = pkgjson.ParseFloat32(value)
		case "position":
			position = pkgjson.ParseVector3(value)
		case "scale":
			scale = pkgjson.ParseVector3(value)
		case "rotation":
			rotation = pkgjson.ParseRotationQuaternion(value)
		case "radius":
			radius = pkgjson.ParseFloat32(value)
		}
	}

	typename := typeNode.(string)

	var l light.Light

	switch typename {
	case "Directional":
		l = light.NewDirectional()
	case "Disk":
		l = light.NewDisk(radius)
	case "Sun":
		l = light.NewDisk(/*0.5 * 0.00935*/0.01)
	case "Point":
		l = light.NewPoint()
	case "Sphere":
		l = light.NewSphere(radius)
	}

	if l == nil {
		return
	}

	l.SetColor(color)
	l.SetLumen(lumen)
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

	c.Init(loader.scene, loader.resourceManager)
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
	rotation := math.MakeIdentityQuaternion()

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

	prop := loader.scene.CreateStaticProp()
	prop.Shape = shape
	prop.Material = material
	prop.SetTransformation(position, scale, rotation)
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