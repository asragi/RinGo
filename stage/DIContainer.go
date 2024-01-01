package stage

type DependencyInjectionContainer struct {
	ValidateAction       ValidateActionFunc
	CalcSkillGrowth      CalcSkillGrowthFunc
	CalcGrowthApply      GrowthApplyFunc
	CalcEarnedItem       CalcEarnedItemFunc
	CalcConsumedItem     CalcConsumedItemFunc
	CalcTotalItem        CalcTotalItemFunc
	StaminaReduction     StaminaReductionFunc
	GetPostActionArgs    GetPostActionArgsFunc
	MakeStageUserExplore createCompensateMakeUserExploreFunc
	MakeUserExplore      makeUserExploreArrayFunc
	GetAllStage          getAllStageFunc
}

func CreateDIContainer() DependencyInjectionContainer {
	validateActionArg := createValidateAction(checkIsExplorePossible)
	return DependencyInjectionContainer{
		ValidateAction:    validateActionArg,
		CalcSkillGrowth:   calcSkillGrowthService,
		CalcGrowthApply:   calcApplySkillGrowth,
		CalcEarnedItem:    calcEarnedItem,
		CalcConsumedItem:  calcConsumedItem,
		CalcTotalItem:     calcTotalItem,
		StaminaReduction:  calcStaminaReduction,
		GetPostActionArgs: GetPostActionArgs,
		MakeUserExplore:   makeUserExplore,
		GetAllStage:       getAllStage,
	}
}
