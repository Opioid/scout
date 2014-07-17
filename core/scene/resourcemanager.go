package scene

type ResourceManager struct {
	MaterialProvider *MaterialProvider
}

func NewResourceManager() *ResourceManager {
	manager := new(ResourceManager)

	manager.MaterialProvider = &MaterialProvider{}

	return manager
}

