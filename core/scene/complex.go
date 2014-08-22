package scene

import (
	"github.com/Opioid/scout/core/resource"
)

type Complex interface {
	Init(scene *Scene, resourceManager *resource.Manager)
}

type ComplexFactory interface {
	New() Complex
}

type ComplexProvider struct {
	factories map[string]ComplexFactory
}

func (p *ComplexProvider) Init() {
	p.factories = make(map[string]ComplexFactory)
}

func (p *ComplexProvider) Register(name string, factory ComplexFactory) {
	p.factories[name] = factory
}

func (p *ComplexProvider) NewComplex(typename string) Complex {
	if f, ok := p.factories[typename]; ok {
		return f.New()
	}

	return nil
}