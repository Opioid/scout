package resource

import (
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/scene/material"
	"github.com/Opioid/scout/core/rendering/texture"
)

type Manager struct {
	shapeProvider *shape.Provider 
	textureProvider *texture.Provider
	materialProvider *material.Provider
}

func NewManager() *Manager {
	manager := new(Manager)

	manager.shapeProvider = &shape.Provider{}
	manager.materialProvider = &material.Provider{}

	return manager
}

func (m *Manager) LoadShape(filename string) shape.Shape {
	return m.shapeProvider.Load(filename)
}

func (m *Manager) LoadTexture2D(filename string, config texture.Config) *texture.Texture2D {
	return m.textureProvider.Load2D(filename, config)
}

func (m *Manager) LoadMaterial(filename string) material.Material {
	return m.materialProvider.Load(filename, m.textureProvider)
}