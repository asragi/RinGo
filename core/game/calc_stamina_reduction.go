package game

import "github.com/asragi/RinGo/core"

type CalcStaminaReductionFunc func(core.StaminaCost, StaminaReducibleRate, []*UserSkillRes) core.StaminaCost

func CalcStaminaReduction(
	baseStamina core.StaminaCost,
	reducibleRate StaminaReducibleRate,
	reductionSkills []*UserSkillRes,
) core.StaminaCost {
	skillLvs := func(skills []*UserSkillRes) []core.SkillLv {
		result := make([]core.SkillLv, len(skills))
		for i, v := range skills {
			result[i] = v.SkillExp.CalcLv()
		}
		return result
	}(reductionSkills)
	skillRate := func(skillLvs []core.SkillLv) float64 {
		result := 1.0
		for _, v := range skillLvs {
			result = v.ApplySkillRate(result)
		}
		return result
	}(skillLvs)
	stamina := ApplyReduction(baseStamina, skillRate, reducibleRate)
	return stamina
}
