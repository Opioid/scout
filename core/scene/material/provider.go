package material

import (
	"github.com/Opioid/scout/core/rendering/material/glass"
	"github.com/Opioid/scout/core/rendering/material/substitute"
	"github.com/Opioid/scout/core/rendering/texture"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"github.com/Opioid/scout/base/math"
	"io/ioutil"
	"encoding/json"
	"fmt"
)


type Provider struct {
	materials map[string]Material

	glassStack 		*glass.BinnedStack
	substituteStack	*substitute.BinnedStack
}

func NewProvider(numWorkers uint32) *Provider {
	p := Provider{}
	p.materials = make(map[string]Material)

	p.glassStack 	  = glass.NewBinnedStack(numWorkers)
	p.substituteStack = substitute.NewBinnedStack(numWorkers)

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

	var material Material

	for key, value := range renderingNode {
		switch key {
		case "Glass":
			material = p.loadGlass(value, tp)
		case "Substitute":
			material = p.loadSubstitute(value, tp)				
		}
	}	

	if material != nil {
		p.materials[filename] = material
	} else {
		fmt.Printf("Load material \"%v\": Couldn't find known material type.\n", filename)
	}

	return material
}

func (p *Provider) loadGlass(i interface{}, tp *texture.Provider) Material {
	node, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	color := math.MakeVector3(0.75, 0.75, 0.75)
	var normalMap *texture.Texture2D

	for key, value := range node {
		switch key {
		case "textures":
			textures, ok := value.([]interface{})

			if !ok {
				break 
			}

			for _, t := range textures {
				texturename, usage := readFilename(t)

				if usage == "Normal" {
					normalMap = tp.Load2D(texturename, texture.Config{Usage: texture.Normal})
				}
			}

		case "color":
			color = pkgjson.ParseVector3(value)
		}
	}

	var material Material

	if normalMap != nil {
		material = glass.NewColorConstant_NormalMap(color, normalMap, p.glassStack)
	} else {		
		material = glass.NewColorConstant(color, p.glassStack)
	}

	return material
}

func (p *Provider) loadSubstitute(i interface{}, tp *texture.Provider) Material {
	node, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	color     := math.MakeVector3(0.75, 0.75, 0.75)
	roughness := float32(1.0)
	metallic  := float32(0.0)
	var colorMap, normalMap *texture.Texture2D

	for key, value := range node {
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
				} else if usage == "Normal" {
					normalMap = tp.Load2D(texturename, texture.Config{Usage: texture.Normal})
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
			material = substitute.NewColorMap_NormalMap(roughness, metallic, colorMap, normalMap, p.substituteStack)
		} else {
			material = substitute.NewColorMap(roughness, metallic, colorMap, p.substituteStack)
		}
	} else {
		if normalMap != nil {
			material = substitute.NewColorConstant_NormalMap(color, roughness, metallic, normalMap, p.substituteStack)
		} else {		
			material = substitute.NewColorConstant(color, roughness, metallic, p.substituteStack)
		}
	}

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