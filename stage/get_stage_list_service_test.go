package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

func TestCreateGetStageListService(t *testing.T) {
	type testRequest struct {
		UserId core.UserId
		Token  core.AccessToken
	}
	type testCase struct {
		request testRequest
		expect  getStageListRes
	}

	createService := CreateGetStageListService(
		stageMasterRepo,
		userStageRepo,
		itemStorageRepo,
		exploreMasterRepo,
		userExploreRepo,
		userSkillRepo,
		consumingItemRepo,
		requiredSkillRepo,
	)

	getStageListService := createService.GetAllStage

	testCases := []testCase{
		{
			request: testRequest{
				UserId: MockUserId,
			},
			expect: getStageListRes{
				Information: []stageInformation{
					{
						StageId: mockStageIds[0],
						IsKnown: true,
						UserExplores: []userExplore{
							{
								ExploreId:  mockStageExploreIds[0],
								IsKnown:    true,
								IsPossible: true,
							},
							{
								ExploreId:  mockStageExploreIds[1],
								IsKnown:    true,
								IsPossible: false,
							},
						},
					},
					{
						StageId: mockStageIds[1],
						IsKnown: true,
					},
				},
			},
		},
	}

	for _, v := range testCases {
		req := v.request
		res, _ := getStageListService(req.UserId, req.Token)
		infos := res.Information
		checkInt(t, "check response length", len(v.expect.Information), len(infos))
		for j, w := range v.expect.Information {
			info := infos[j]
			check(t, string(w.StageId), string(info.StageId))
			checkInt(t, "check response explore length", len(w.UserExplores), len(info.UserExplores))
			for k, x := range w.UserExplores {
				explore := info.UserExplores[k]
				check(t, string(x.ExploreId), string(explore.ExploreId))
				checkBool(t, "IsKnown", bool(x.IsKnown), bool(explore.IsKnown))
				checkBool(t, "IsPossible", bool(x.IsPossible), bool(explore.IsPossible))
			}
		}
	}
}
