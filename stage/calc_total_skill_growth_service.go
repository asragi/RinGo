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

type growthApplyFunc func(core.UserId, core.AccessToken, []skillGrowthResult) []growthApplyResult

type growthApplyRes struct {
	Create growthApplyFunc
}

func calcSkillGrowthApplyResultService(
	userSkillRepo UserSkillRepo,
) growthApplyRes {
	create := func(
		userId core.UserId,
		token core.AccessToken,
		skillGrowth []skillGrowthResult,
	) []growthApplyResult {
		toSkillId := func(skillGrowthResults []skillGrowthResult) []core.SkillId {
			result := make([]core.SkillId, len(skillGrowthResults))
			for i, v := range skillGrowthResults {
				result[i] = v.SkillId
			}
			return result
		}

		makeSkillGrowthMap := func(skillGrowthResults []skillGrowthResult) map[core.SkillId]skillGrowthResult {
			result := make(map[core.SkillId]skillGrowthResult)
			for _, v := range skillGrowthResults {
				result[v.SkillId] = v
			}
			return result
		}

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

		skillGrowthMap := makeSkillGrowthMap(skillGrowth)
		skillsRes, err := userSkillRepo.BatchGet(userId, toSkillId(skillGrowth), token)
		if err != nil {
			return []growthApplyResult{}
		}

		result := make([]growthApplyResult, len(skillsRes.Skills))
		for i, v := range skillsRes.Skills {
			userSkill := v
			result[i] = applySkillGrowth(userSkill, skillGrowthMap[userSkill.SkillId])
		}
		return result
	}

	return growthApplyRes{
		Create: create,
	}
}
