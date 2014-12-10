package material

import (
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/core/rendering/texture"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"github.com/Opioid/scout/base/math"
	"io/ioutil"
	"encoding/json"
	"fmt"
)


type Provider struct {

}

func (p *Provider) Load(filename string, tp *texture.Provider) Material {
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

	if colorMap != nil {
		if normalMap != nil {
			return material.NewSubstitute_ColorMap_NormalMap(roughness, metallic, colorMap, normalMap)
		} else {
			return material.NewSubstitute_ColorMap(roughness, metallic, colorMap)
		}
	} else {
		if normalMap != nil {
			return material.NewSubstitute_ColorConstant_NormalMap(color, roughness, metallic, normalMap)
		} else {		
			return material.NewSubstitute_ColorConstant(color, roughness, metallic)
		}
	}
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