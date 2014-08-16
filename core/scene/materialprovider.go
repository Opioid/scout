package scene

import (
	"github.com/Opioid/scout/core/rendering/material"
	"github.com/Opioid/scout/core/rendering/texture"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"github.com/Opioid/scout/base/math"
	"io/ioutil"
	"encoding/json"
	"fmt"
)


type MaterialProvider struct {

}

func (p *MaterialProvider) Load(filename string, m *ResourceManager) Material {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil
	}

	var document interface{}
	if err = json.Unmarshal(data, &document); err != nil {
		fmt.Println(err)
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

//	material := &Material{Color: math.Vector3{0.75, 0.75, 0.75}, Roughness: 0.9 }

//	material := new(material.Substitute_ColorOnly)

	var color math.Vector3
	var roughness float32
	var colorMap texture.Sampler2D

	for key, value := range renderingNode {
		switch key {
		case "textures":
			textures, ok := value.([]interface{})

			if !ok {
				break 
			}

			for _, t := range textures {
				filename = readFilename(t)
				if colorTexture := m.LoadTexture2D(filename); colorTexture != nil {
					colorMap = texture.NewSampler_nearest(colorTexture, new(texture.AddressMode_repeat))
				}
			}

		case "color":
			color = pkgjson.ParseVector3(value)
		case "roughness":
			roughness = float32(value.(float64))
		}
	}

	if colorMap != nil {
		return material.NewSubstitute_ColorMap(color, roughness, colorMap)
	} else {
		return material.NewSubstitute_ColorConstant(color, roughness)
	}
}

func readFilename(i interface{}) string {
	node, ok := i.(map[string]interface{})

	if !ok {
		return ""
	}

	filename := node["file"].(string)

	return filename
}