package scene

import (
	pkgjson "github.com/Opioid/scout/base/parsing/json"
	"github.com/Opioid/scout/base/math"
	"io/ioutil"
	"encoding/json"
	"fmt"
)


type MaterialProvider struct {

}

func (p *MaterialProvider) Load(filename string) *Material {
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

	material := &Material{Color: math.Vector3{0.75, 0.75, 0.75}, Roughness: 0.9 }

	for key, value := range renderingNode {
		switch key {
		case "color":
			material.Color = pkgjson.ParseVector3(value)
		case "roughness":
			material.Roughness = float32(value.(float64))
		}
	}

	return material
}