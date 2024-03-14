package stage

// Deprecated: avoid using DependencyInjectionContainer
type DependencyInjectionContainer struct {
	GetPostActionArgs                GetPostActionArgsFunc
	MakeStageUserExplore             CreateCompensateMakeUserExploreFunc
	MakeUserExplore                  MakeUserExploreArrayFunc
	GetAllStage                      GetAllStageFunc
	CreateGetUserResourceServiceFunc CreateGetUserResourceServiceFunc
}

func CreateDIContainer() DependencyInjectionContainer {
	return DependencyInjectionContainer{
		GetPostActionArgs:                GetPostActionArgs,
		MakeUserExplore:                  MakeUserExplore,
		GetAllStage:                      getAllStage,
		CreateGetUserResourceServiceFunc: CreateGetUserResourceService,
		MakeStageUserExplore:             CompensateMakeUserExplore,
	}
}
