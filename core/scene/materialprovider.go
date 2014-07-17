package scene

type MaterialProvider struct {

}

func (p *MaterialProvider) Load(name string) *Material {
	material := new(Material)

	return material
}