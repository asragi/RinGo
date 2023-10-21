package stage

import (
	"github.com/asragi/RinGo/core"
)

type skillGrowthResult struct {
	SkillId core.SkillId
	GainSum GainingPoint
}

type CalcSkillGrowthFunc func([]SkillGrowthData, int) []skillGrowthResult

func calcSkillGrowthService(execCount int, gainingData []SkillGrowthData) []skillGrowthResult {
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
