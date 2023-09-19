package stage

import (
	"testing"

	"github.com/asragi/RinGo/core"
)

var (
	itemMasterRepo      = CreateMockItemMasterRepo()
	itemStorageRepo     = CreateMockItemStorageRepo()
	userExploreRepo     = createMockUserExploreRepo()
	conditionRepo       = createMockExploreConditionRepo()
	exploreMasterRepo   = createMockExploreMasterRepo()
	skillMasterRepo     = createMockSkillMasterRepo()
	userSkillRepo       = createMockUserSkillRepo()
	userStageRepo       = createMockUserStageRepo()
	stageMasterRepo     = createMockStageMasterRepo()
	skillGrowthDataRepo = createMockSkillGrowthDataRepo()
)

func check(t *testing.T, expect string, actual string) {
	if expect != actual {
		t.Errorf("want %s, actual %s", expect, actual)
	}
}

func checkBool(t *testing.T, title string, expect bool, actual bool) {
	if expect != actual {
		t.Errorf("%s: want %t, actual %t", title, expect, actual)
	}
}

func checkInt(t *testing.T, title string, expect int, actual int) {
	if expect != actual {
		t.Errorf("%s: want %d, actual %d", title, expect, actual)
	}
}

func TestCreateItemService(t *testing.T) {
	type testRequest struct {
		userId core.UserId
		itemId core.ItemId
	}

	type testExplore struct {
		exploreId  ExploreId
		name       core.DisplayName
		isKnown    core.IsKnown
		isPossible core.IsPossible
	}

	type testExpect struct {
		price    core.Price
		stock    core.Stock
		explores []testExplore
	}

	type testCase struct {
		request testRequest
		expect  testExpect
	}

	itemService := CreateItemService(
		itemMasterRepo,
		itemStorageRepo,
		exploreMasterRepo,
		userExploreRepo,
		skillMasterRepo,
		userSkillRepo,
		conditionRepo)
	getUserItemDetail := itemService.GetUserItemDetail

	testCases := []testCase{
		{
			request: testRequest{
				itemId: MockItems[0].ItemId,
				userId: MockUserId,
			},
			expect: testExpect{
				price: MockItems[0].Price,
				stock: 20,
				explores: []testExplore{
					{
						exploreId:  mockExploreIds[0],
						name:       mockExploreMaster[MockItems[0].ItemId][0].DisplayName,
						isKnown:    true,
						isPossible: true,
					},
					{
						exploreId:  mockExploreIds[1],
						name:       mockExploreMaster[MockItems[0].ItemId][1].DisplayName,
						isKnown:    false,
						isPossible: false,
					},
				},
			},
		},
	}
	// test
	for _, v := range testCases {
		targetId := v.request.itemId
		req := GetUserItemDetailReq{
			UserId: v.request.userId,
			ItemId: targetId,
		}
		res := getUserItemDetail(req)
		// check proper id
		if res.ItemId != targetId {
			t.Errorf("want %s, actual %s", targetId, res.ItemId)
		}

		// check proper master data
		expect := v.expect
		if res.Price != expect.price {
			t.Errorf("want %d, actual %d", expect.price, res.Price)
		}

		// check proper user storage data
		targetStock := expect.stock
		if res.Stock != targetStock {
			t.Errorf("want %d, actual %d", targetStock, res.Stock)
		}

		// check explore
		if len(res.UserExplores) != len(expect.explores) {
			t.Errorf("want %d, actual %d", len(expect.explores), len(res.UserExplores))
		}
		for j, w := range expect.explores {
			actual := res.UserExplores[j]
			if w.exploreId != actual.ExploreId {
				t.Errorf("want %s, actual %s", w.exploreId, actual.ExploreId)
			}
			check(t, string(w.name), string(actual.DisplayName))
			checkBool(t, "isKnown", bool(w.isKnown), bool(actual.IsKnown))
			checkBool(t, "isPossible", bool(w.isPossible), bool(actual.IsPossible))
		}
	}
}

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
		conditionRepo,
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
						UserExplores: []UserExplore{
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
		res := getStageListService(req.UserId, req.Token)
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

func TestCalcSkillGrowthService(t *testing.T) {
	type testRequest struct {
		exploreId ExploreId
		execCount int
	}
	type testCase struct {
		request testRequest
		expect  []skillGrowthResult
	}

	testCases := []testCase{
		{
			request: testRequest{
				exploreId: mockExploreIds[0],
				execCount: 3,
			},
			expect: []skillGrowthResult{
				{
					SkillId: mockSkillIds[0],
					GainSum: 30,
				},
				{
					SkillId: mockSkillIds[1],
					GainSum: 30,
				},
			},
		},
	}

	service := createCalcSkillGrowthService(skillGrowthDataRepo)

	for _, v := range testCases {
		req := v.request
		res := service.Calc(req.exploreId, req.execCount)
		checkInt(t, "skill growth response length", len(v.expect), len(res))
		for i, w := range v.expect {
			result := res[i]
			check(t, string(w.SkillId), string(result.SkillId))
			checkInt(t, "check skill exp gain sum", int(w.GainSum), int(result.GainSum))
		}
	}
}
