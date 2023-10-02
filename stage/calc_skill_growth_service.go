package stage

import "github.com/asragi/RinGo/core"

type skillGrowthResult struct {
	SkillId core.SkillId
	GainSum GainingPoint
}

type calcSkillGrowthFunc func(ExploreId, int) []skillGrowthResult

type calcSkillGrowthService struct {
	Calc calcSkillGrowthFunc
}

func createCalcSkillGrowthService(
	skillGrowthDataRepo SkillGrowthDataRepo,
) calcSkillGrowthService {
	calcSkillGrowth := func(
		exploreId ExploreId,
		execCount int,
	) []skillGrowthResult {
		calcSkillGrowth := func(execCount int, gainingData []SkillGrowthData) []skillGrowthResult {
			growth := make([]skillGrowthResult, len(gainingData))
			for i := range gainingData {
				data := gainingData[i]
				growth[i] = skillGrowthResult{
					SkillId: data.SkillId,
					GainSum: data.GainingPoint.Multiply(execCount),
				}
			}
			return growth
		}
		skillGrowthList := skillGrowthDataRepo.BatchGet(exploreId)
		result := calcSkillGrowth(execCount, skillGrowthList)
		return result
	}
	return calcSkillGrowthService{Calc: calcSkillGrowth}
}
