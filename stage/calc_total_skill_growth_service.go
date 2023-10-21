package stage

import "github.com/asragi/RinGo/core"

type growthApplyResult struct {
	SkillId   core.SkillId
	GainSum   GainingPoint
	BeforeLv  core.SkillLv
	BeforeExp core.SkillExp
	AfterLv   core.SkillLv
	AfterExp  core.SkillExp
	WasLvUp   bool
}

type GrowthApplyFunc func([]UserSkillRes, []skillGrowthResult) []growthApplyResult

type growthApplyRes struct {
	Create GrowthApplyFunc
}

func calcApplySkillGrowth(userSkills []UserSkillRes, skillGrowth []skillGrowthResult) []growthApplyResult {
	applySkillGrowth := func(userSkill UserSkillRes, skillGrowth skillGrowthResult) growthApplyResult {
		if userSkill.SkillId != skillGrowth.SkillId {
			// TODO: proper error handling
			panic("invalid apply skill growth!")
		}
		beforeExp := userSkill.SkillExp
		afterExp := skillGrowth.GainSum.ApplyTo(beforeExp)
		beforeLv := beforeExp.CalcLv()
		afterLv := afterExp.CalcLv()
		wasLvUp := beforeLv != afterLv
		return growthApplyResult{
			SkillId:   userSkill.SkillId,
			GainSum:   skillGrowth.GainSum,
			BeforeLv:  beforeLv,
			BeforeExp: beforeExp,
			AfterLv:   afterLv,
			AfterExp:  afterExp,
			WasLvUp:   wasLvUp,
		}
	}
	makeSkillGrowthMap := func(skillGrowthResults []skillGrowthResult) map[core.SkillId]skillGrowthResult {
		result := make(map[core.SkillId]skillGrowthResult)
		for _, v := range skillGrowthResults {
			result[v.SkillId] = v
		}
		return result
	}

	skillGrowthMap := makeSkillGrowthMap(skillGrowth)

	result := make([]growthApplyResult, len(userSkills))
	for i, v := range userSkills {
		userSkill := v
		result[i] = applySkillGrowth(userSkill, skillGrowthMap[userSkill.SkillId])
	}
	return result
}
