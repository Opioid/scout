package scene

import (
	"github.com/Opioid/scout/core/rendering/surrounding"
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/scene/material"
	"github.com/Opioid/scout/core/resource"
	"github.com/Opioid/scout/base/math"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"github.com/Opioid/scout/base/file"
	"io/ioutil"
	"encoding/json"
	"os"
	"path/filepath"
	"fmt"
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

	var scene jsonScene
	if err = json.Unmarshal(data, &scene); err != nil {
		return err
	}	 

	loader.loadSurrounding(&scene.Surrounding)

	loader.loadEntities(scene.Entities)

	loader.scene.Compile()

	return nil
}

func (loader *Loader) loadSurrounding(js *jsonSurrounding) {
	switch js.Class {
	case "Uniform":
		color := math.MakeVector3FromArray(js.Color)  
		loader.scene.Surrounding = surrounding.NewUniform(color)
	case "Textured":
		usaCache := false

		if usaCache {
			filenameBase := file.WithoutExt(filepath.Base(js.Texture.File))

			diffuseTextureName  := filenameBase  + "_diffuse.sui"
			specularTextureName := filenameBase  + "_specular.sui"

			diffuseTexture  := loadCachedTexture(diffuseTextureName)
			specularTexture := loadCachedTexture(specularTextureName)

			if diffuseTexture != nil && specularTexture != nil {
				loader.scene.Surrounding = surrounding.NewSphereFromCache(diffuseTexture, specularTexture)
				fmt.Println("Found cached surrounding.")
			} else {
				if sphericalTexture := loader.resourceManager.LoadTexture2D(js.Texture.File, texture.Config{Usage: texture.RGBA}); 
				   sphericalTexture != nil {
					s := surrounding.NewSphere(sphericalTexture)

					loader.scene.Surrounding = s

					saveCachedTexture(diffuseTextureName, s.DiffuseTexture())
					saveCachedTexture(specularTextureName, s.SpecularTexture())

					fmt.Println("Created cached surrounding.")
				} 
			}
		} else {
			if sphericalTexture := loader.resourceManager.LoadTexture2D(js.Texture.File, texture.Config{Usage: texture.RGBA}); 
			   sphericalTexture != nil {
				s := surrounding.NewSphere(sphericalTexture)

				loader.scene.Surrounding = s
			}
		}		
	}

	if loader.scene.Surrounding == nil {
		loader.scene.Surrounding = surrounding.NewUniform(math.MakeVector3(1.0, 0.5, 0.75))
	}	
}

func (loader *Loader) loadEntities(entities []jsonEntity) {
	for e := range entities {
		switch entities[e].Type {
		case "Light":
			loader.loadLight(&entities[e])
		case "Prop":
			loader.loadProp(&entities[e])
	//	case "Complex":
	//		loader.loadComplex(entityNode)
		}
	}
}

type jsonScene struct {
	Surrounding jsonSurrounding

	Entities []jsonEntity
}

type jsonSurrounding struct {
	Class string

	Texture jsonTexture

	Color [3]float32
}

type jsonTexture struct {
	File string
}

type jsonEntity struct {
	Type, Class  string

	Position [3]float32
	Scale    [3]float32
	Rotation [3]float32

	Shape jsonShape

	Keyframes []jsonKeyframe

	Materials []string

	Color [3]float32
	Lumen float32
	Radius float32
}

type jsonShape struct {
	Class, File string
}

type jsonKeyframe struct {
	Position [3]float32
	Scale    [3]float32
	Rotation [3]float32
}

func (loader *Loader) loadLight(e *jsonEntity) {
//	shape := loader.loadShape(&e.Shape)

//	materials := loader.loadMaterials(e.Materials)

	var l light.Light

	switch e.Class {
	case "Directional":
		l = light.NewDirectional()
	case "Disk":
		l = light.NewDisk(e.Radius)
	case "Sun":
		l = light.NewDisk(0.0125)
	case "Point":
		l = light.NewPoint()
	case "Sphere":
		l = light.NewSphere(e.Radius)
	}

	if l == nil {
		return
	}

	position := math.MakeVector3FromArray(e.Position)
	scale := math.MakeVector3FromArray(e.Scale)
	rotation := pkgjson.MakeRotationQuaternion(e.Rotation)

	color := math.MakeVector3FromArray(e.Color)

	l.SetColor(color)
	l.SetLumen(e.Lumen)
	l.Prop().SetTransformation(position, scale, rotation)
//	l.Entity().Animation = animation

	loader.scene.AddLight(l)


/*
	prop := loader.scene.CreateProp()
	prop.Shape = shape
	prop.Materials = materials
	prop.SetTransformation(position, scale, rotation)	
	*/	
}

/*
func (loader *Loader) loadLight(i interface{}) {
	lightNode, ok := i.(map[string]interface{})

	if !ok {
		return
	}

	typeNode, ok := lightNode["type"]

	if !ok {
		return
	}

	var position math.Vector3
	scale := math.MakeIdentityVector3()
	rotation := math.MakeIdentityQuaternion()
	var animation entity.Animation

	color := math.MakeVector3(1.0, 1.0, 1.0)
	lumen := float32(1.0)	
	radius := float32(1.0)

	for key, value := range lightNode {
		switch key {
		case "position":
			position = pkgjson.ParseVector3(value)
		case "scale":
			scale = pkgjson.ParseVector3(value)
		case "rotation":
			rotation = pkgjson.ParseRotationQuaternion(value)
		case "keyframes":
			animation = entity.MakeAnimationFromJson(value)						
		case "color":
			color = pkgjson.ParseVector3(value)
		case "lumen":
			lumen = pkgjson.ParseFloat32(value)
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
		l = light.NewDisk(0.0125)
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
	l.Entity().SetTransformation(position, scale, rotation)
	l.Entity().Animation = animation

	loader.scene.AddLight(l)
}
/*
/*
func (loader *Loader) loadActor(i interface{}) {
	actorNode, ok := i.(map[string]interface{})

	if !ok {
		return
	}

	s, ok := actorNode["shape"]

	if !ok {
		return
	}

	shape := loader.loadShape(s)

	if shape == nil {
		return
	}

	m, ok := actorNode["materials"]

	if !ok {
		return
	}

	materials := loader.loadMaterials(m)

	if materials == nil {
		return
	}

	var position math.Vector3
	scale := math.MakeIdentityVector3()
	rotation := math.MakeIdentityQuaternion()
	var animation entity.Animation	

	for key, value := range actorNode {
		switch key {
		case "position":
			position = pkgjson.ParseVector3(value)
		case "scale":
			scale = pkgjson.ParseVector3(value)
		case "rotation":
			rotation = pkgjson.ParseRotationQuaternion(value)
		case "keyframes":
			animation = entity.MakeAnimationFromJson(value)					
		}
	}

	a := loader.scene.CreateActor()
	a.Shape = shape
	a.Materials = materials
	a.SetTransformation(position, scale, rotation)
	a.Animation = animation	
}
*/
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

func (loader *Loader) loadProp(e *jsonEntity) {
	shape := loader.loadShape(&e.Shape)
	if shape == nil {
		return
	}

	materials := loader.loadMaterials(e.Materials)
	if materials == nil {
		return
	}

	position := math.MakeVector3FromArray(e.Position)
	scale := math.MakeVector3FromArray(e.Scale)
	rotation := pkgjson.MakeRotationQuaternion(e.Rotation)

	prop := loader.scene.CreateProp()
	prop.Shape = shape
	prop.Materials = materials
	prop.SetTransformation(position, scale, rotation)	
}

/*
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

	m, ok := propNode["materials"]

	if !ok {
		return
	}

	materials := loader.loadMaterials(m)

	if materials == nil {
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
	prop.Materials = materials
	prop.SetTransformation(position, scale, rotation)
}
*/

func (loader *Loader) loadShape(s *jsonShape) shape.Shape {
	if len(s.Class) > 0 {
		switch s.Class {
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
	} else if len(s.File) > 0 {
		return loader.resourceManager.LoadShape(s.File)
	}

	return nil
}

func (loader *Loader) loadMaterials(materialNames []string) []material.Material {
	materials := make([]material.Material, len(materialNames))

	for i, m := range materialNames {
		material := loader.resourceManager.LoadMaterial(m)

		if material == nil {
			return nil
		}

		materials[i] = material
	}

	return materials
}	

func saveCachedTexture(filename string, t *texture.Texture2D) error {
	fo, err := os.Create("../cache/" + filename)
	defer fo.Close()

	if err == nil {
		if err := texture.Save(fo, t); err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func loadCachedTexture(filename string) *texture.Texture2D {
	fi, err := os.Open("../cache/" + filename)
	defer fi.Close()

	if err == nil {
		texture, err := texture.Load(fi)

		if err != nil {
			fmt.Println(err)
			return nil
		} else {
			return texture
		}
	}	

	return nil
}