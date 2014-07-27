package scene

import (
	"github.com/Opioid/scout/core/scene/shape"
)

type ResourceManager struct {
	shapeProvider *shape.Provider 
	materialProvider *MaterialProvider
}

func NewResourceManager() *ResourceManager {
	manager := new(ResourceManager)

	manager.shapeProvider = &shape.Provider{}
	manager.materialProvider = &MaterialProvider{}

	return manager
}

func (m *ResourceManager) LoadShape(filename string) shape.Shape {
	return m.shapeProvider.Load(filename)
}

func (m *ResourceManager) LoadMaterial(filename string) *Material {
	return m.materialProvider.Load(filename)
}