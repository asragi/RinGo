package game

import "github.com/asragi/RinGo/core"

type CalcStaminaReductionFunc func(core.Stamina, StaminaReducibleRate, []*UserSkillRes) core.Stamina

func CalcStaminaReduction(
	baseStamina core.Stamina,
	reducibleRate StaminaReducibleRate,
	reductionSkills []*UserSkillRes,
) core.Stamina {
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
