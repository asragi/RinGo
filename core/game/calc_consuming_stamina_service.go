package game

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type ExploreStaminaPair struct {
	ExploreId      ActionId
	ReducedStamina core.StaminaCost
}

type CalcConsumingStaminaFunc func(
	context.Context,
	core.UserId,
	[]ActionId,
) ([]*ExploreStaminaPair, error)

func CreateCalcConsumingStaminaService(
	fetchUserSkills FetchUserSkillFunc,
	fetchExploreMaster FetchExploreMasterFunc,
	fetchReductionSkills FetchReductionStaminaSkillFunc,
) CalcConsumingStaminaFunc {
	return func(ctx context.Context, userId core.UserId, exploreIds []ActionId) (
		[]*ExploreStaminaPair,
		error,
	) {
		handleError := func(err error) ([]*ExploreStaminaPair, error) {
			return nil, fmt.Errorf("error on batch calc mockReducedStamina: %w", err)
		}
		explores, err := fetchExploreMaster(ctx, exploreIds)
		if err != nil {
			return handleError(err)
		}
		exploreMap := func(explores []*GetExploreMasterRes) map[ActionId]GetExploreMasterRes {
			result := map[ActionId]GetExploreMasterRes{}
			for _, v := range explores {
				result[v.ExploreId] = *v
			}
			return result
		}(explores)
		reductionStaminaSkills, err := fetchReductionSkills(ctx, exploreIds)
		if err != nil {
			return handleError(err)
		}
		reductionSkillMap := func(reductionSkills []*StaminaReductionSkillPair) map[ActionId][]core.SkillId {
			result := map[ActionId][]core.SkillId{}
			for _, v := range reductionSkills {
				result[v.ExploreId] = append(result[v.ExploreId], v.SkillId)
			}
			return result
		}(reductionStaminaSkills)

		allRequiredSkill := func(skills []*StaminaReductionSkillPair) []core.SkillId {
			check := map[core.SkillId]bool{}
			var result []core.SkillId
			for _, v := range skills {
				skillId := v.SkillId
				if _, ok := check[skillId]; ok {
					continue
				}
				check[skillId] = true
				result = append(result, skillId)
			}
			return result
		}(reductionStaminaSkills)

		allSkills, err := fetchUserSkills(ctx, userId, allRequiredSkill)
		if err != nil {
			return handleError(err)
		}

		allSkillsMap := func(
			userId core.UserId,
			skills []*UserSkillRes,
			skillId []core.SkillId,
		) map[core.SkillId]*UserSkillRes {
			result := map[core.SkillId]*UserSkillRes{}
			for _, v := range skills {
				result[v.SkillId] = v
			}
			for _, v := range skillId {
				if _, ok := result[v]; !ok {
					result[v] = &UserSkillRes{
						UserId:   userId,
						SkillId:  v,
						SkillExp: 0,
					}
				}
			}
			return result
		}(userId, allSkills.Skills, allRequiredSkill)

		reductionSkillResMap := func(
			allSkillsMap map[core.SkillId]*UserSkillRes,
			reductionSkills map[ActionId][]core.SkillId,
		) map[ActionId][]*UserSkillRes {
			result := map[ActionId][]*UserSkillRes{}
			for k, v := range reductionSkills {
				for _, w := range v {
					if _, ok := result[k]; !ok {
						result[k] = []*UserSkillRes{}
					}
					result[k] = append(result[k], allSkillsMap[w])
				}
			}
			return result
		}(allSkillsMap, reductionSkillMap)

		result := func(
			idArr []ActionId,
			exploreMap map[ActionId]GetExploreMasterRes,
			reductionSkillMap map[ActionId][]*UserSkillRes,
		) []*ExploreStaminaPair {
			result := make([]*ExploreStaminaPair, len(exploreMap))
			index := 0
			for _, v := range idArr {
				explore := exploreMap[v]
				baseStamina := explore.ConsumingStamina
				reducibleRate := explore.StaminaReducibleRate
				stamina := CalcStaminaReduction(baseStamina, reducibleRate, reductionSkillMap[v])
				result[index] = &ExploreStaminaPair{
					ExploreId:      v,
					ReducedStamina: stamina,
				}
				index++
			}
			return result
		}(exploreIds, exploreMap, reductionSkillResMap)
		return result, nil
	}
}
