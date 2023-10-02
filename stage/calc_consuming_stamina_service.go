package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
)

type calcConsumingStaminaFunc func(core.UserId, core.AccessToken, ExploreId) (core.Stamina, error)

type createCalcConsumingStaminaServiceRes struct {
	Calc calcConsumingStaminaFunc
}

func createCalcConsumingStaminaService(
	userSkillRepo UserSkillRepo,
	exploreMasterRepo ExploreMasterRepo,
	reductionSkillRepo ReductionStaminaSkillRepo,
) createCalcConsumingStaminaServiceRes {
	calc := func(userId core.UserId, token core.AccessToken, exploreId ExploreId) (core.Stamina, error) {
		handleError := func(err error) (core.Stamina, error) {
			return 0, fmt.Errorf("error on calculating stamina: %w", err)
		}
		explore, err := exploreMasterRepo.Get(exploreId)
		if err != nil {
			return handleError(err)
		}
		baseStamina := explore.ConsumingStamina
		reducibleRate := explore.StaminaReducibleRate
		reductionSkills, err := reductionSkillRepo.Get(exploreId)
		if err != nil {
			return handleError(err)
		}
		userSkillsRes, err := userSkillRepo.BatchGet(userId, reductionSkills, token)
		if err != nil {
			return handleError(err)
		}
		skillLvs := func(skills []UserSkillRes) []core.SkillLv {
			result := make([]core.SkillLv, len(skills))
			for i, v := range skills {
				result[i] = v.SkillExp.CalcLv()
			}
			return result
		}(userSkillsRes.Skills)
		skillRate := func(skillLvs []core.SkillLv) float64 {
			result := 1.0
			for _, v := range skillLvs {
				result = v.ApplySkillRate(result)
			}
			return result
		}(skillLvs)
		stamina := ApplyReduction(baseStamina, skillRate, reducibleRate)
		return stamina, nil
	}

	return createCalcConsumingStaminaServiceRes{
		Calc: calc,
	}
}
