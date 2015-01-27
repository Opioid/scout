package scene

import (
	"github.com/Opioid/scout/core/rendering/surrounding"
	"github.com/Opioid/scout/core/rendering/texture"
	lightmaterial "github.com/Opioid/scout/core/rendering/material/light"
	"github.com/Opioid/scout/core/scene/prop"
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
	switch js.Type {
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
		default:
			loader.loadComplex(&entities[e])
		}
	}
}

type jsonScene struct {
	Surrounding jsonSurrounding

	Entities []jsonEntity
}

type jsonSurrounding struct {
	Type string

	Texture jsonTexture

	Color [3]float32
}

type jsonTexture struct {
	File string
}

type jsonEntity struct {
	Type  string
	Shape jsonShape

	Position [3]float32
	Scale    [3]float32
	Rotation [3]float32

	Keyframes []jsonKeyframe

	Materials []string

	Color [3]float32
	Lumen float32
	Radius float32
}

type jsonShape struct {
	Type, File string
}

type jsonKeyframe struct {
	Position [3]float32
	Scale    [3]float32
	Rotation [3]float32
}

func (loader *Loader) loadLight(e *jsonEntity) {
	position := math.MakeVector3FromArray(e.Position)
	scale := math.MakeVector3FromArray(e.Scale)
	rotation := pkgjson.MakeRotationQuaternion(e.Rotation)

	color := math.MakeVector3FromArray(e.Color)

	var l light.Light

	switch e.Shape.Type {
	case "Directional":
		l = light.NewDirectional()
	case "Disk":
		l = light.NewDisk()
	case "Sun":
		l = light.NewDisk()
		scale = math.MakeVector3FromScalar(0.0125)
	case "Point":
		l = light.NewPoint()
	case "Sphere":
		l = light.NewSphere(loader.shape(e.Shape.Type))
	}

	if l == nil && len(e.Shape.File) > 0 {
		// It probably is a mesh light
		l = light.NewMesh(loader.resourceManager.LoadShape(e.Shape.File))
	}

	if l == nil {
		// We didn't find any light
		return
	}

	l.SetColor(color)
	l.SetLumen(e.Lumen)

	if l.Prop().Shape != nil {
		materials := make([]material.Material, 1)
		materials[0] = lightmaterial.NewColorConstant(l)
		l.Prop().Materials = materials
	}

	
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
func (loader *Loader) loadComplex(e *jsonEntity) {
	c := loader.scene.CreateComplex(e.Type)
	if c == nil {
		return
	}

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

	prop :=	prop.NewProp()
	prop.Shape = shape
	prop.Materials = materials
	prop.SetTransformation(position, scale, rotation)
	loader.scene.AddProp(prop)
}

func (loader *Loader) loadShape(s *jsonShape) shape.Shape {
	if len(s.Type) > 0 {
		return loader.shape(s.Type)
	} else if len(s.File) > 0 {
		return loader.resourceManager.LoadShape(s.File)
	}

	return nil
}

func (loader *Loader) shape(typename string) shape.Shape {
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

	return nil
}

func (loader *Loader) loadMaterials(materialNames []string) []material.Material {
	if len(materialNames) == 0 {
		return nil
	}

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