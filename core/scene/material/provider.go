package material

import (
	pkgmaterial "github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/core/rendering/texture"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"github.com/Opioid/scout/base/math"
	"io/ioutil"
	"encoding/json"
	"fmt"
)


type Provider struct {
	materials map[string]Material
}

func NewProvider() *Provider {
	p := Provider{}
	p.materials = make(map[string]Material)
	return &p
}

func (p *Provider) Load(filename string, tp *texture.Provider) Material {
	if material, ok := p.materials[filename]; ok {
		return material
	}

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil
	}

	var document interface{}
	if err = json.Unmarshal(data, &document); err != nil {
		fmt.Printf("Load material \"%v\": %v\n", filename, err)
		return nil
	}

	root := document.(map[string]interface{})

	r, ok := root["rendering"]

	if !ok {
		return nil
	}

	renderingNode, ok := r.(map[string]interface{})

	if !ok {
		return nil
	}

	color     := math.MakeVector3(0.75, 0.75, 0.75)
	roughness := float32(1)
	metallic  := float32(0)
	var colorMap, normalMap *texture.Texture2D

	for key, value := range renderingNode {
		switch key {
		case "textures":
			textures, ok := value.([]interface{})

			if !ok {
				break 
			}

			for _, t := range textures {
				texturename, usage := readFilename(t)

				if usage == "Color" {
					colorMap = tp.Load2D(texturename, texture.Config{Usage: texture.RGBA})
				} else if usage == "Normals" {
					normalMap = tp.Load2D(texturename, texture.Config{Usage: texture.Normals})
				}

			}

		case "color":
			color = pkgjson.ParseVector3(value)
		case "roughness":
			roughness = float32(value.(float64))
		case "metallic":
			metallic = float32(value.(float64))
		}
	}

	var material Material

	if colorMap != nil {
		if normalMap != nil {
			material = pkgmaterial.NewSubstitute_ColorMap_NormalMap(roughness, metallic, colorMap, normalMap)
		} else {
			material = pkgmaterial.NewSubstitute_ColorMap(roughness, metallic, colorMap)
		}
	} else {
		if normalMap != nil {
			material = pkgmaterial.NewSubstitute_ColorConstant_NormalMap(color, roughness, metallic, normalMap)
		} else {		
			material = pkgmaterial.NewSubstitute_ColorConstant(color, roughness, metallic)
		}
	}

	p.materials[filename] = material

	return material
}

func readFilename(i interface{}) (string, string) {
	node, ok := i.(map[string]interface{})

	if !ok {
		return "", ""
	}

	filename := pkgjson.ReadString(node, "file", "")
	usage    := pkgjson.ReadString(node, "usage", "Color")

	return filename, usage
}