package scene

import (
	"github.com/Opioid/scout/core/rendering/material"
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"github.com/Opioid/scout/base/math"
	"io/ioutil"
	"encoding/json"
	"fmt"
)


type MaterialProvider struct {

}

func (p *MaterialProvider) Load(filename string) Material {
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

	for key, value := range renderingNode {
		switch key {
		case "color":
			color = pkgjson.ParseVector3(value)
		case "roughness":
			roughness = float32(value.(float64))
		}
	}

	return material.NewSubstitute_ColorOnly(color, roughness)
}