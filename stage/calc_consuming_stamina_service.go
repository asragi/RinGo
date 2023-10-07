package stage

import (
	"fmt"

	"github.com/asragi/RinGo/core"
)

type ExploreStaminaPair struct {
	ExploreId      ExploreId
	ReducedStamina core.Stamina
}

type calcConsumingStaminaFunc func(core.UserId, core.AccessToken, ExploreId) (core.Stamina, error)
type calcBatchConsumingStaminaFunc func(core.UserId, core.AccessToken, []GetExploreMasterRes) ([]ExploreStaminaPair, error)

type createCalcConsumingStaminaServiceRes struct {
	// Deprecated: Replace with BatchCalc()
	Calc      calcConsumingStaminaFunc
	BatchCalc calcBatchConsumingStaminaFunc
}

func CreateCalcConsumingStaminaService(
	userSkillRepo UserSkillRepo,
	exploreMasterRepo ExploreMasterRepo,
	reductionSkillRepo ReductionStaminaSkillRepo,
) createCalcConsumingStaminaServiceRes {
	staminaReduction := func(
		baseStamina core.Stamina,
		reducibleRate StaminaReducibleRate,
		reductionSkills []UserSkillRes) core.Stamina {
		skillLvs := func(skills []UserSkillRes) []core.SkillLv {
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
		stamina := staminaReduction(baseStamina, reducibleRate, userSkillsRes.Skills)
		return stamina, nil
	}

	batchCalc := func(userId core.UserId, token core.AccessToken, explores []GetExploreMasterRes) ([]ExploreStaminaPair, error) {
		handleError := func(err error) ([]ExploreStaminaPair, error) {
			return nil, fmt.Errorf("error on batch calc stamina: %w", err)
		}
		exploreMap := func(explores []GetExploreMasterRes) map[ExploreId]GetExploreMasterRes {
			result := map[ExploreId]GetExploreMasterRes{}
			for _, v := range explores {
				result[v.ExploreId] = v
			}
			return result
		}(explores)
		exploreIds := func(objects []GetExploreMasterRes) []ExploreId {
			result := make([]ExploreId, len(objects))
			for i, v := range objects {
				result[i] = v.ExploreId
			}
			return result
		}(explores)
		reductionStaminaSkills, err := reductionSkillRepo.BatchGet(exploreIds)
		if err != nil {
			return handleError(err)
		}
		reductionSkillMap := func(reductionSkills []BatchGetReductionStaminaSkill) map[ExploreId][]core.SkillId {
			result := map[ExploreId][]core.SkillId{}
			for _, v := range reductionSkills {
				result[v.ExploreId] = v.Skills
			}
			return result
		}(reductionStaminaSkills)

		allRequiredSkill := func(skills []BatchGetReductionStaminaSkill) []core.SkillId {
			check := map[core.SkillId]bool{}
			result := []core.SkillId{}
			for _, v := range skills {
				for _, w := range v.Skills {
					if _, ok := check[w]; ok {
						continue
					}
					check[w] = true
					result = append(result, w)
				}
			}
			return result
		}(reductionStaminaSkills)

		allSkills, err := userSkillRepo.BatchGet(userId, allRequiredSkill, token)
		if err != nil {
			return handleError(err)
		}

		allSkillsMap := func(skills []UserSkillRes) map[core.SkillId]UserSkillRes {
			result := map[core.SkillId]UserSkillRes{}
			for _, v := range skills {
				result[v.SkillId] = v
			}
			return result
		}(allSkills.Skills)

		reductionSkillResMap := func(
			allSkillsMap map[core.SkillId]UserSkillRes,
			reductionSkills map[ExploreId][]core.SkillId,
		) map[ExploreId][]UserSkillRes {
			result := map[ExploreId][]UserSkillRes{}
			for k, v := range reductionSkills {
				for _, w := range v {
					if _, ok := result[k]; !ok {
						result[k] = []UserSkillRes{}
					}
					result[k] = append(result[k], allSkillsMap[w])
				}
			}
			return result
		}(allSkillsMap, reductionSkillMap)

		result := func(
			exploreMap map[ExploreId]GetExploreMasterRes,
			reductionSkillMap map[ExploreId][]UserSkillRes,
		) []ExploreStaminaPair {
			result := make([]ExploreStaminaPair, len(exploreMap))
			index := 0
			for k, v := range exploreMap {
				explore := v
				baseStamina := explore.ConsumingStamina
				reducibleRate := explore.StaminaReducibleRate
				stamina := staminaReduction(baseStamina, reducibleRate, reductionSkillMap[k])
				result[index] = ExploreStaminaPair{
					ExploreId:      k,
					ReducedStamina: stamina,
				}
				index++
			}
			return result
		}(exploreMap, reductionSkillResMap)
		return result, nil
	}

	return createCalcConsumingStaminaServiceRes{
		Calc:      calc,
		BatchCalc: batchCalc,
	}
}
