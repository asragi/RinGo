package explore

// Deprecated: avoid using DependencyInjectionContainer
type DependencyInjectionContainer struct {
	GetAllStage                      GetAllStageFunc
	CreateGetUserResourceServiceFunc CreateGetUserResourceServiceFunc
}

func CreateDIContainer() DependencyInjectionContainer {
	return DependencyInjectionContainer{
		GetAllStage:                      GetAllStage,
		CreateGetUserResourceServiceFunc: CreateGetUserResourceService,
	}
}
