package shape

import (
	"io/ioutil"
	"encoding/json"
)

type Provider struct {

}

func (p *Provider) Load(filename string) Shape {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil
	}

	var document interface{}
	if err = json.Unmarshal(data, &document); err != nil {
		return nil
	}

	root := document.(map[string]interface{})

	t, hasType := root["type"]

	if !hasType {
		return nil
	}

	switch t {
	case "Sphere":
		return new(Sphere)
	case "Plane":
		return new(Plane)
	default:
		return nil
	}

	return nil
}
