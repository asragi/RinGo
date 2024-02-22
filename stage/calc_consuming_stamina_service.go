package stage

import (
	"fmt"
	"github.com/asragi/RinGo/core"
)

type StaminaReductionFunc func(core.Stamina, StaminaReducibleRate, []UserSkillRes) core.Stamina

func calcStaminaReduction(
	baseStamina core.Stamina,
	reducibleRate StaminaReducibleRate,
	reductionSkills []UserSkillRes,
) core.Stamina {
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

type ExploreStaminaPair struct {
	ExploreId      ExploreId
	ReducedStamina core.Stamina
}

type CalcBatchConsumingStaminaFunc func(
	core.UserId,
	[]ExploreId,
) (
	[]ExploreStaminaPair,
	error,
)

func CreateCalcConsumingStaminaService(
	fetchUserSkills FetchUserSkillFunc,
	fetchExploreMaster FetchExploreMasterFunc,
	fetchReductionSkills FetchReductionStaminaSkillFunc,
) CalcBatchConsumingStaminaFunc {
	batchCalc := func(userId core.UserId, exploreIds []ExploreId) (
		[]ExploreStaminaPair,
		error,
	) {
		handleError := func(err error) ([]ExploreStaminaPair, error) {
			return nil, fmt.Errorf("error on batch calc mockReducedStamina: %w", err)
		}
		explores, err := fetchExploreMaster(exploreIds)
		if err != nil {
			return handleError(err)
		}
		exploreMap := func(explores []GetExploreMasterRes) map[ExploreId]GetExploreMasterRes {
			result := map[ExploreId]GetExploreMasterRes{}
			for _, v := range explores {
				result[v.ExploreId] = v
			}
			return result
		}(explores)
		reductionStaminaSkills, err := fetchReductionSkills(exploreIds)
		if err != nil {
			return handleError(err)
		}
		reductionSkillMap := func(reductionSkills []BatchGetReductionStaminaSkill) map[ExploreId][]core.SkillId {
			result := map[ExploreId][]core.SkillId{}
			for _, v := range reductionSkills {
				skillIds := func(pair []StaminaReductionSkillPair) []core.SkillId {
					skillIdArr := make([]core.SkillId, len(pair))
					for i, v := range pair {
						skillIdArr[i] = v.SkillId
					}
					return skillIdArr
				}(v.Skills)
				result[v.ExploreId] = skillIds
			}
			return result
		}(reductionStaminaSkills)

		allRequiredSkill := func(skills []BatchGetReductionStaminaSkill) []core.SkillId {
			check := map[core.SkillId]bool{}
			var result []core.SkillId
			for _, v := range skills {
				for _, w := range v.Skills {
					skillId := w.SkillId
					if _, ok := check[skillId]; ok {
						continue
					}
					check[skillId] = true
					result = append(result, skillId)
				}
			}
			return result
		}(reductionStaminaSkills)

		allSkills, err := fetchUserSkills(userId, allRequiredSkill)
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
				stamina := calcStaminaReduction(baseStamina, reducibleRate, reductionSkillMap[k])
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

	return batchCalc
}
