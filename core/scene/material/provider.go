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
	var colorSampler texture.Sampler2D

	for key, value := range renderingNode {
		switch key {
		case "textures":
			textures, ok := value.([]interface{})

			if !ok {
				break 
			}

			for _, t := range textures {
				texturename, _ := readFilename(t)
				if colorTexture := tp.Load2D(texturename, false); colorTexture != nil {
					colorSampler = texture.NewSampler2D_linear(colorTexture, new(texture.AddressMode_repeat))
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

	if colorSampler != nil {
		return material.NewSubstitute_ColorMap(color, roughness, metallic, colorSampler)
	} else {
		return material.NewSubstitute_ColorConstant(color, roughness, metallic)
	}
}

func readFilename(i interface{}) (string, string) {
	node, ok := i.(map[string]interface{})

	if !ok {
		return "", ""
	}

	filename := pkgjson.ReadString(node, "file", "")
	usage    := pkgjson.ReadString(node, "usage", "")

	return filename, usage
}