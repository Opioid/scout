package scene

import (
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/rendering/texture"
)

type ResourceManager struct {
	shapeProvider *shape.Provider 
	textureProvider *texture.Provider
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

func (m *ResourceManager) LoadTexture2D(filename string) *texture.Texture2D {
	return m.textureProvider.Load2D(filename)
}

func (m *ResourceManager) LoadMaterial(filename string) Material {
	return m.materialProvider.Load(filename, m)
}