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
	userSkillMap := func(userSkills []UserSkillRes) map[core.SkillId]UserSkillRes {
		result := make(map[core.SkillId]UserSkillRes)
		for _, v := range userSkills {
			result[v.SkillId] = v
		}
		return result
	}(userSkills)

	result := make([]growthApplyResult, len(skillGrowth))
	for i, v := range skillGrowth {
		userSkill, ok := userSkillMap[v.SkillId]
		if !ok {
			userSkill = UserSkillRes{
				SkillId:  v.SkillId,
				SkillExp: 0,
			}
		}
		result[i] = applySkillGrowth(userSkill, v)
	}
	return result
}
