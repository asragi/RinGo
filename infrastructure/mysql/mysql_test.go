package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/core/game/shelf/reservation"
	"github.com/asragi/RinGo/infrastructure"
	"github.com/asragi/RinGo/location"
	"github.com/asragi/RinGo/test"
	"github.com/asragi/RinGo/utils"
	"testing"
	"time"
)

type userTest struct {
	UserId             core.UserId                `db:"user_id"`
	Name               core.UserName              `db:"name"`
	MaxStamina         core.MaxStamina            `db:"max_stamina"`
	Fund               core.Fund                  `db:"fund"`
	StaminaRecoverTime time.Time                  `db:"stamina_recover_time"`
	HashedPassword     auth.HashedPassword        `db:"hashed_password"`
	Popularity         reservation.ShopPopularity `db:"popularity"`
}

var TestCompleted = errors.New("test completed")

type ApplyUserTestOption func(*userTest)

func createTestUser(options ...ApplyUserTestOption) *userTest {
	user := userTest{
		UserId:             "created-test-user",
		Name:               "test-name",
		MaxStamina:         6000,
		Fund:               100000,
		StaminaRecoverTime: test.MockTime(),
		HashedPassword:     "test-password",
		Popularity:         reservation.ShopPopularity(50),
	}
	for _, option := range options {
		option(&user)
	}
	return &user
}

func TestCreateCheckUserExistence(t *testing.T) {
	type testCase struct {
		userId      core.UserId
		expectedErr error
	}

	ctx := test.MockCreateContext()
	errorUserId := core.UserId("error-user")
	testUser := createTestUser(func(user *userTest) { user.UserId = errorUserId })
	_, err := dba.Exec(
		ctx,
		insertTestUserQuery,
		testUser,
	)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}
	defer func() {
		_, err := dba.Exec(ctx, "DELETE FROM ringo.users WHERE user_id = :user_id", testUser)
		if err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}
	}()
	testCases := []testCase{
		{userId: "valid-user", expectedErr: nil},
		{userId: errorUserId, expectedErr: auth.UserAlreadyExistsError},
	}
	for _, v := range testCases {
		checkUserExistence := CreateCheckUserExistence(dba.Query)
		testErr := checkUserExistence(ctx, v.userId)
		if !errors.Is(testErr, v.expectedErr) {
			t.Errorf("got: %v, expect: %v", errors.Unwrap(testErr), v.expectedErr)
		}
	}
}

func TestCreateGetUserPassword(t *testing.T) {
	type testCase struct {
		userId         core.UserId
		hashedPassword auth.HashedPassword
	}

	testCases := []testCase{
		{userId: "test-password-user", hashedPassword: "test-password"},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		user := createTestUser(func(user *userTest) { user.HashedPassword = v.hashedPassword; user.UserId = v.userId })
		_, err := dba.Exec(
			ctx,
			insertTestUserQuery,
			user,
		)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}
		getUserPassword := CreateGetUserPassword(dba.Query)
		res, err := getUserPassword(ctx, v.userId)
		if err != nil {
			t.Errorf("failed to fetch user password: %v", err)
		}
		if res != v.hashedPassword {
			t.Errorf("got: %v, expect: %v", res, v.hashedPassword)
		}
		func() {
			_, err := dba.Exec(ctx, "DELETE FROM ringo.users WHERE user_id = :user_id", user)
			if err != nil {
				t.Fatalf("failed to delete user: %v", err)
			}
		}()
	}
}

func TestCreateInsertNewUser(t *testing.T) {
	type testCase struct {
		UserId             core.UserId                `db:"user_id"`
		Name               core.UserName              `db:"name"`
		HashedPassword     auth.HashedPassword        `db:"hashed_password"`
		InitialFund        core.Fund                  `db:"fund"`
		InitialStamina     core.MaxStamina            `db:"max_stamina"`
		InitialPopularity  reservation.ShopPopularity `db:"popularity"`
		StaminaRecoverTime time.Time                  `db:"stamina_recover_time"`
	}

	testCases := []testCase{
		{
			UserId:             "test-insert-user",
			Name:               "test-name",
			HashedPassword:     "test-password",
			InitialFund:        3456,
			InitialStamina:     5678,
			InitialPopularity:  50,
			StaminaRecoverTime: test.MockTime().In(location.UTC()),
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		insertNewUser := CreateInsertNewUser(
			dba.Exec,
			v.InitialFund,
			v.InitialStamina,
			v.InitialPopularity,
			test.MockTime,
		)
		txErr := dba.Transaction(
			ctx, func(ctx context.Context) error {
				err := insertNewUser(ctx, v.UserId, v.Name, v.HashedPassword)
				if err != nil {
					t.Fatalf("failed to insert new user: %v", err)
					return err
				}
				rows, err := dba.Query(
					ctx,
					"SELECT user_id, name, hashed_password, fund, max_stamina, stamina_recover_time FROM ringo.users WHERE user_id = :user_id",
					&v,
				)
				if err != nil {
					t.Fatalf("failed to fetch user: %v", err)
					return err
				}
				if !rows.Next() {
					t.Fatalf("failed to fetch user: %v", err)
					return err
				}
				var res testCase
				err = rows.StructScan(&res)
				if err != nil {
					t.Fatalf("failed to fetch user: %v", err)
					return err
				}
				err = rows.Close()
				if err != nil {
					t.Fatalf("failed to close rows: %v", err)
					return err
				}
				// DeepEqual doesn't work well with time.Time due to location setting
				if res.UserId != v.UserId {
					t.Errorf("got: %v, expect: %v", res.UserId, v.UserId)
				}
				if res.Name != v.Name {
					t.Errorf("got: %v, expect: %v", res.Name, v.Name)
				}
				if res.HashedPassword != v.HashedPassword {
					t.Errorf("got: %v, expect: %v", res.HashedPassword, v.HashedPassword)
				}
				if res.InitialFund != v.InitialFund {
					t.Errorf("got: %v, expect: %v", res.InitialFund, v.InitialFund)
				}
				if !test.DeepEqual(res.InitialStamina, v.InitialStamina) {
					t.Errorf("got: %v, expect: %v", res.InitialStamina, v.InitialStamina)
				}
				if !res.StaminaRecoverTime.Equal(v.StaminaRecoverTime) {
					t.Errorf("got: %v, expect: %v", res.StaminaRecoverTime, v.StaminaRecoverTime)
				}
				_, err = dba.Exec(ctx, "DELETE FROM ringo.users WHERE user_id = :user_id", &v)
				if err != nil {
					t.Fatalf("failed to delete user: %v", err)
					return err
				}
				return TestCompleted
			},
		)
		if !errors.Is(txErr, TestCompleted) {
			t.Fatalf("failed: %v", txErr)
		}
	}
}

func TestCreateGetResourceMySQL(t *testing.T) {
	type testCase struct {
		UserId             core.UserId     `db:"user_id"`
		MaxStamina         core.MaxStamina `db:"max_stamina"`
		StaminaRecoverTime time.Time       `db:"stamina_recover_time"`
		Fund               core.Fund       `db:"fund"`
	}

	testCases := []testCase{
		{
			UserId:             "test-get-resource-user",
			MaxStamina:         6000,
			StaminaRecoverTime: test.MockTime(),
			Fund:               100000,
		},
	}

	for _, v := range testCases {
		user := createTestUser(
			func(user *userTest) {
				user.UserId = v.UserId
				user.MaxStamina = v.MaxStamina
				user.StaminaRecoverTime = v.StaminaRecoverTime
				user.Fund = v.Fund
			},
		)
		expectedRes := &game.GetResourceRes{
			UserId:             v.UserId,
			MaxStamina:         v.MaxStamina,
			StaminaRecoverTime: core.StaminaRecoverTime(v.StaminaRecoverTime),
			Fund:               v.Fund,
		}
		ctx := test.MockCreateContext()
		_, err := dba.Exec(
			ctx,
			insertTestUserQuery,
			user,
		)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}
		fetchResource := CreateGetResourceMySQL(dba.Query)
		res, err := fetchResource(ctx, v.UserId)
		if err != nil {
			t.Errorf("failed to fetch resource: %v", err)
		}
		if test.DeepEqual(res, expectedRes) {
			t.Errorf("got: %+v, expect: %+v", res, expectedRes)
			if res.UserId != expectedRes.UserId {
				t.Errorf("got: %v, expect: %v", res.UserId, expectedRes.UserId)
			}
			if res.MaxStamina != expectedRes.MaxStamina {
				t.Errorf("got: %v, expect: %v", res.MaxStamina, expectedRes.MaxStamina)
			}
			if res.StaminaRecoverTime != expectedRes.StaminaRecoverTime {
				t.Errorf("got: %v, expect: %v", res.StaminaRecoverTime, expectedRes.StaminaRecoverTime)
			}
		}
		_, err = dba.Exec(ctx, "DELETE FROM ringo.users WHERE user_id = :user_id", user)
		if err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}
	}
}

func TestCreateUpdateStamina(t *testing.T) {
	type testCase struct {
		UserId             core.UserId
		StaminaRecoverTime time.Time
		AfterRecoverTime   time.Time
	}

	testCases := []testCase{
		{
			UserId:             "test-update-stamina-user",
			StaminaRecoverTime: test.MockTime(),
			AfterRecoverTime:   test.MockTime().Add(time.Hour).In(location.UTC()),
		},
	}

	for _, v := range testCases {
		user := createTestUser(
			func(user *userTest) {
				user.UserId = v.UserId
				user.StaminaRecoverTime = v.StaminaRecoverTime
			},
		)
		ctx := test.MockCreateContext()
		_, err := dba.Exec(
			ctx,
			insertTestUserQuery,
			user,
		)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}
		updateStamina := CreateUpdateStamina(dba.Exec)
		err = updateStamina(ctx, v.UserId, core.StaminaRecoverTime(v.AfterRecoverTime))
		if err != nil {
			t.Fatalf("failed to update stamina: %v", err)
		}
		rows, err := dba.Query(
			ctx,
			"SELECT user_id, stamina_recover_time FROM ringo.users WHERE user_id = :user_id",
			user,
		)
		if err != nil {
			t.Fatalf("failed to fetch user: %v", err)
		}
		if !rows.Next() {
			t.Fatalf("failed to fetch user: %v", err)
		}
		var res userTest
		err = rows.StructScan(&res)
		if err != nil {
			t.Fatalf("failed to fetch user: %v", err)
		}
		if !res.StaminaRecoverTime.Equal(v.AfterRecoverTime) {
			t.Errorf("got: %v, expect: %v", res.StaminaRecoverTime, v.AfterRecoverTime)
		}
		_, err = dba.Exec(ctx, "DELETE FROM ringo.users WHERE user_id = :user_id", user)
		if err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}
	}
}

func TestCreateGetItemMasterMySQL(t *testing.T) {
	type testCase struct {
		itemId []core.ItemId
	}

	testCases := []testCase{
		{itemId: []core.ItemId{"1"}},
		{itemId: []core.ItemId{"1", "2"}},
	}

	for _, v := range testCases {
		fetchItemMaster := CreateGetItemMasterMySQL(dba.Query)
		ctx := test.MockCreateContext()
		res, err := fetchItemMaster(ctx, v.itemId)
		if err != nil {
			t.Errorf("failed to fetch item master: %v", err)
		}
		if len(res) != len(v.itemId) {
			t.Errorf("got: %d, expect: %d", len(res), len(v.itemId))
		}
	}
}

func TestCreateGetStageMaster(t *testing.T) {
	type testCase struct {
		stageId []explore.StageId
	}

	testCases := []testCase{
		{stageId: []explore.StageId{"1"}},
		{stageId: []explore.StageId{"1", "2"}},
	}

	for _, v := range testCases {
		fetchStageMaster := CreateGetStageMaster(dba.Query)
		ctx := test.MockCreateContext()
		res, err := fetchStageMaster(ctx, v.stageId)
		if err != nil {
			t.Errorf("failed to fetch stage master: %v", err)
		}
		if len(res) != len(v.stageId) {
			t.Errorf("got: %d, expect: %d", len(res), len(v.stageId))
		}
	}
}

func TestCreateGetAllStageMaster(t *testing.T) {
	fetchAllStageMaster := CreateGetAllStageMaster(dba.Query)
	ctx := test.MockCreateContext()
	res, err := fetchAllStageMaster(ctx)
	if err != nil {
		t.Errorf("failed to fetch all stage master: %v", err)
	}
	if len(res) == 0 {
		t.Errorf("got: %d, expect: >0", len(res))
	}
}

func TestCreateGetExploreMasterMySQL(t *testing.T) {
	type testCase struct {
		exploreIds []game.ExploreId
	}

	testCases := []testCase{
		{exploreIds: []game.ExploreId{"1"}},
		{exploreIds: []game.ExploreId{"1", "2"}},
	}

	for _, v := range testCases {
		fetchExploreMaster := CreateGetExploreMasterMySQL(dba.Query)
		ctx := test.MockCreateContext()
		res, err := fetchExploreMaster(ctx, v.exploreIds)
		if err != nil {
			t.Errorf("failed to fetch explore master: %v", err)
		}
		if len(res) != len(v.exploreIds) {
			t.Errorf("got: %d, expect: %d", len(res), len(v.exploreIds))
		}
	}
}

func TestCreateGetSkillMaster(t *testing.T) {
	type testCase struct {
		skillIds []core.SkillId
	}

	testCases := []testCase{
		{skillIds: []core.SkillId{"1"}},
		{skillIds: []core.SkillId{"1", "2"}},
	}

	for _, v := range testCases {
		fetchSkillMaster := CreateGetSkillMaster(dba.Query)
		ctx := test.MockCreateContext()
		res, err := fetchSkillMaster(ctx, v.skillIds)
		if err != nil {
			t.Fatalf("failed to fetch skill master: %v", err)
		}
		if len(res) != len(v.skillIds) {
			t.Errorf("got: %d, expect: %d", len(res), len(v.skillIds))
		}
	}
}

func TestCreateGetEarningItem(t *testing.T) {
	type testCase struct {
		exploreId   game.ExploreId
		expectedRes []*game.EarningItem
	}

	testCases := []testCase{
		{
			exploreId: "1",
			expectedRes: []*game.EarningItem{
				{
					ItemId:      "1",
					MinCount:    50,
					MaxCount:    100,
					Probability: 1,
				},
			},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		fetchEarningItem := CreateGetEarningItem(dba.Query)
		res, err := fetchEarningItem(ctx, v.exploreId)
		if err != nil {
			t.Fatalf("failed to fetch earning item: %v", err)
		}
		if !test.DeepEqual(res, v.expectedRes) {
			t.Errorf("got: %+v, expect: %+v", res, v.expectedRes)
			for i := range res {
				r := res[i]
				e := v.expectedRes[i]
				if r.ItemId != e.ItemId {
					t.Errorf("got: %v, expect: %v", r.ItemId, e.ItemId)
				}
				if r.MinCount != e.MinCount {
					t.Errorf("got: %v, expect: %v", r.MinCount, e.MinCount)
				}
				if r.MaxCount != e.MaxCount {
					t.Errorf("got: %v, expect: %v", r.MaxCount, e.MaxCount)
				}
				if r.Probability != e.Probability {
					t.Errorf("got: %v, expect: %v", r.Probability, e.Probability)
				}
			}
		}
	}
}

func TestCreateGetConsumingItem(t *testing.T) {
	type testCase struct {
		exploreIds  []game.ExploreId
		expectedRes []*game.ConsumingItem
	}

	testCases := []testCase{
		{
			exploreIds: []game.ExploreId{"2"},
			expectedRes: []*game.ConsumingItem{
				{
					ExploreId:       "2",
					ItemId:          "1",
					MaxCount:        1,
					ConsumptionProb: 1,
				},
			},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		fetchConsumingItem := CreateGetConsumingItem(dba.Query)
		res, err := fetchConsumingItem(ctx, v.exploreIds)
		if err != nil {
			t.Fatalf("failed to fetch consuming item: %v", err)
		}
		if !test.DeepEqual(res, v.expectedRes) {
			t.Errorf("got: %+v, expect: %+v", res, v.expectedRes)
			for i := range res {
				r := res[i]
				e := v.expectedRes[i]
				if r.ExploreId != e.ExploreId {
					t.Errorf("got: %v, expect: %v", r.ExploreId, e.ExploreId)
				}
				if r.ItemId != e.ItemId {
					t.Errorf("got: %v, expect: %v", r.ItemId, e.ItemId)
				}
				if r.MaxCount != e.MaxCount {
					t.Errorf("got: %v, expect: %v", r.MaxCount, e.MaxCount)
				}
				if r.ConsumptionProb != e.ConsumptionProb {
					t.Errorf("got: %v, expect: %v", r.ConsumptionProb, e.ConsumptionProb)
				}
			}
		}
	}
}

func TestCreateGetRequiredSkills(t *testing.T) {
	type testCase struct {
		exploreIds  []game.ExploreId
		expectedRes []*game.RequiredSkill
	}

	testCases := []testCase{
		{
			exploreIds: []game.ExploreId{"5"},
			expectedRes: []*game.RequiredSkill{
				{
					ExploreId:  "5",
					SkillId:    "20",
					RequiredLv: 10,
				},
			},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		fetchRequiredSkills := CreateGetRequiredSkills(dba.Query)
		res, err := fetchRequiredSkills(ctx, v.exploreIds)
		if err != nil {
			t.Fatalf("failed to fetch required skills: %v", err)
		}
		if !test.DeepEqual(res, v.expectedRes) {
			t.Errorf("got: %+v, expect: %+v", res, v.expectedRes)
			for i := range res {
				r := res[i]
				e := v.expectedRes[i]
				if r.ExploreId != e.ExploreId {
					t.Errorf("got: %v, expect: %v", r.ExploreId, e.ExploreId)
				}
				if r.SkillId != e.SkillId {
					t.Errorf("got: %v, expect: %v", r.SkillId, e.SkillId)
				}
			}
		}
	}
}

func TestCreateGetSkillGrowth(t *testing.T) {
	type testCase struct {
		exploreId   game.ExploreId
		expectedRes []*game.SkillGrowthData
	}

	testCases := []testCase{
		{
			exploreId: "1",
			expectedRes: []*game.SkillGrowthData{
				{
					ExploreId:    "1",
					SkillId:      "1",
					GainingPoint: 10,
				},
				{
					ExploreId:    "1",
					SkillId:      "2",
					GainingPoint: 10,
				},
				{
					ExploreId:    "1",
					SkillId:      "14",
					GainingPoint: 10,
				},
			},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		fetchSkillGrowth := CreateGetSkillGrowth(dba.Query)
		res, err := fetchSkillGrowth(ctx, v.exploreId)
		if err != nil {
			t.Fatalf("failed to fetch skill growth: %v", err)
		}
		if !test.DeepEqual(res, v.expectedRes) {
			t.Errorf("got: %+v, expect: %+v", res, v.expectedRes)
			for i := range res {
				r := res[i]
				e := v.expectedRes[i]
				if r.ExploreId != e.ExploreId {
					t.Errorf("got: %v, expect: %v", r.ExploreId, e.ExploreId)
				}
				if r.SkillId != e.SkillId {
					t.Errorf("got: %v, expect: %v", r.SkillId, e.SkillId)
				}
				if r.GainingPoint != e.GainingPoint {
					t.Errorf("got: %v, expect: %v", r.GainingPoint, e.GainingPoint)
				}
			}
		}
	}
}

func TestCreateGetReductionSkill(t *testing.T) {
	type testCase struct {
		exploreId   []game.ExploreId
		expectedRes []*game.StaminaReductionSkillPair
	}

	testCases := []testCase{
		{
			exploreId: []game.ExploreId{"1"},
			expectedRes: []*game.StaminaReductionSkillPair{
				{
					ExploreId: "1",
					SkillId:   "1",
				},
				{
					ExploreId: "1",
					SkillId:   "2",
				},
				{
					ExploreId: "1",
					SkillId:   "14",
				},
			},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		fetchReductionSkill := CreateGetReductionSkill(dba.Query)
		res, err := fetchReductionSkill(ctx, v.exploreId)
		if err != nil {
			t.Fatalf("failed to fetch reduction skill: %v", err)
		}
		if !test.DeepEqual(res, v.expectedRes) {
			t.Errorf("got: %+v, expect: %+v", res, v.expectedRes)
			for i := range res {
				r := res[i]
				e := v.expectedRes[i]
				if r.ExploreId != e.ExploreId {
					t.Errorf("got: %v, expect: %v", r.ExploreId, e.ExploreId)
				}
				if r.SkillId != e.SkillId {
					t.Errorf("got: %v, expect: %v", r.SkillId, e.SkillId)
				}
			}
		}
	}
}

func TestCreateStageExploreRelation(t *testing.T) {
	type testCase struct {
		stageId  []explore.StageId
		expected []*explore.StageExploreIdPairRow
	}

	testCases := []testCase{
		{
			stageId: []explore.StageId{"1"},
			expected: []*explore.StageExploreIdPairRow{
				{
					StageId:   "1",
					ExploreId: "1",
				},
				{
					StageId:   "1",
					ExploreId: "5",
				},
			},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		fetchStageExploreRelation := CreateStageExploreRelation(dba.Query)
		res, err := fetchStageExploreRelation(ctx, v.stageId)
		if err != nil {
			t.Fatalf("failed to fetch stage explore relation: %v", err)
		}
		if !test.DeepEqual(res, v.expected) {
			t.Errorf("got: %+v, expect: %+v", res, v.expected)
			for i := range res {
				r := res[i]
				e := v.expected[i]
				if r.StageId != e.StageId {
					t.Errorf("got: %v, expect: %v", r.StageId, e.StageId)
				}
				if r.ExploreId != e.ExploreId {
					t.Errorf("got: %v, expect: %v", r.ExploreId, e.ExploreId)
				}
			}
		}
	}
}

func TestCreateItemExploreRelation(t *testing.T) {
	type testCase struct {
		itemId   core.ItemId
		expected []game.ExploreId
	}

	testCases := []testCase{
		{
			itemId:   "1",
			expected: []game.ExploreId{"2", "3", "4"},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		fetchItemExploreRelation := CreateItemExploreRelation(dba.Query)
		res, err := fetchItemExploreRelation(ctx, v.itemId)
		if err != nil {
			t.Fatalf("failed to fetch item explore relation: %v", err)
		}
		if !test.DeepEqual(res, v.expected) {
			t.Errorf("got: %+v, expect: %+v", res, v.expected)
			for i := range res {
				r := res[i]
				e := v.expected[i]
				if r != e {
					t.Errorf("got: %v, expect: %v", r, e)
				}
			}
		}
	}
}

func TestCreateGetUserExplore(t *testing.T) {
	type testCase struct {
		userId   core.UserId
		expected []*game.ExploreUserData
	}

	type exploreData struct {
		UserId    core.UserId    `db:"user_id"`
		ExploreId game.ExploreId `db:"explore_id"`
		IsKnown   core.IsKnown   `db:"is_known"`
	}
	testCases := []testCase{
		{
			userId: testUserId,
			expected: []*game.ExploreUserData{
				{
					ExploreId: "1",
					IsKnown:   true,
				},
				{
					ExploreId: "2",
					IsKnown:   false,
				},
			},
		},
	}

	type userIdRes struct {
		UserId core.UserId `db:"user_id"`
	}
	allUserId := func() []*userIdRes {
		var res []*userIdRes
		for _, v := range testCases {
			res = append(res, &userIdRes{UserId: v.userId})
		}
		return res
	}()

	defer func() {
		ctx := test.MockCreateContext()
		userIds := func() []core.UserId {
			var res []core.UserId
			for _, w := range allUserId {
				res = append(res, w.UserId)
			}
			return res
		}()
		allUserIdString := spreadString(infrastructure.UserIdsToString(userIds))
		_, err := dba.Exec(
			ctx,
			fmt.Sprintf(`DELETE FROM ringo.user_explore_data WHERE user_id IN (%s)`, allUserIdString),
			nil,
		)
		if err != nil {
			t.Fatalf("failed to delete user explore: %v", err)
		}
	}()

	for _, v := range testCases {
		data := func(original []*game.ExploreUserData) []*exploreData {
			var res []*exploreData
			for _, w := range original {
				res = append(
					res, &exploreData{
						UserId:    v.userId,
						ExploreId: w.ExploreId,
						IsKnown:   w.IsKnown,
					},
				)
			}
			return res
		}(v.expected)
		exploreIds := func() []game.ExploreId {
			var res []game.ExploreId
			for _, w := range v.expected {
				res = append(res, w.ExploreId)
			}
			return res
		}()
		ctx := test.MockCreateContext()
		_, err := dba.Exec(
			ctx,
			`INSERT INTO ringo.user_explore_data (user_id, explore_id, is_known) VALUES (:user_id, :explore_id, :is_known)`,
			data,
		)
		if err != nil {
			t.Fatalf("failed to insert user explore: %v", err)
		}
		fetchUserExplore := CreateGetUserExplore(dba.Query)
		res, err := fetchUserExplore(ctx, v.userId, exploreIds)
		if err != nil {
			t.Fatalf("failed to fetch user explore: %v", err)
		}
		if !test.DeepEqual(res, v.expected) {
			t.Errorf("got: %+v, expect: %+v", res, v.expected)
		}
	}
}

func TestCreateGetUserStageData(t *testing.T) {
	type testCase struct {
		userId   core.UserId
		stageIds []explore.StageId
		expected []*explore.UserStage
	}

	type stageData struct {
		UserId  core.UserId     `db:"user_id"`
		StageId explore.StageId `db:"stage_id"`
		IsKnown bool            `db:"is_known"`
	}

	stageDataSet := []*stageData{
		{
			UserId:  testUserId,
			StageId: "1",
			IsKnown: true,
		},
		{
			UserId:  testUserId,
			StageId: "2",
			IsKnown: false,
		},
	}

	ctx := test.MockCreateContext()
	// insert test data
	_, err := dba.Exec(
		ctx,
		`INSERT INTO ringo.user_stage_data (user_id, stage_id, is_known) VALUES (:user_id, :stage_id, :is_known)`,
		stageDataSet,
	)
	if err != nil {
		t.Fatalf("failed to insert user stage data: %v", err)
	}
	defer func() {
		_, err := dba.Exec(ctx, `DELETE FROM ringo.user_stage_data WHERE user_id = :user_id`, stageDataSet[0])
		if err != nil {
			t.Fatalf("failed to delete user stage data: %v", err)
		}
	}()

	testCases := []testCase{
		{
			userId:   testUserId,
			stageIds: []explore.StageId{"1", "2"},
			expected: []*explore.UserStage{
				{
					StageId: "1",
					IsKnown: true,
				},
				{
					StageId: "2",
					IsKnown: false,
				},
			},
		},
	}

	for _, v := range testCases {
		fetchUserStageData := CreateGetUserStageData(dba.Query)
		res, err := fetchUserStageData(ctx, v.userId, v.stageIds)
		if err != nil {
			t.Fatalf("failed to fetch user stage data: %v", err)
		}
		if !test.DeepEqual(res, v.expected) {
			t.Errorf("got: %+v, expect: %+v", res, v.expected)
		}
	}
}

func TestCreateUpdateFund(t *testing.T) {
	type testCase struct {
		req []*game.UserFundPair
	}

	testCases := []testCase{
		{
			req: []*game.UserFundPair{
				{
					UserId: "test_fund_user",
					Fund:   10001,
				},
				{
					UserId: "test_fund_user2",
					Fund:   10002,
				},
				{
					UserId: "test_fund_user3",
					Fund:   10003,
				},
			},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		testUsers := make([]*userTest, len(v.req))
		for i, w := range v.req {
			user := createTestUser(func(u *userTest) { u.UserId = w.UserId })
			testUsers[i] = user
			_, err := dba.Exec(
				ctx,
				insertTestUserQuery,
				user,
			)
			if err != nil {
				t.Fatalf("failed to insert user: %v", err)
			}
		}
		updateFund := CreateUpdateFund(dba.Exec)
		err := updateFund(ctx, v.req)
		if err != nil {
			t.Fatalf("failed to update fund: %v", err)
		}
		for j, user := range testUsers {
			rows, err := dba.Query(
				ctx,
				"SELECT fund FROM ringo.users WHERE user_id = :user_id",
				user,
			)
			if err != nil {
				t.Fatalf("failed to fetch user: %v", err)
			}
			if !rows.Next() {
				t.Fatalf("failed to fetch user: %v", err)
			}
			var res core.Fund
			err = rows.Scan(&res)
			if err != nil {
				t.Fatalf("failed to fetch user: %v", err)
			}
			expect := v.req[j].Fund
			if res != expect {
				t.Errorf("got: %v, expect: %v", res, expect)
			}
			_, err = dba.Exec(ctx, "DELETE FROM ringo.users WHERE user_id = :user_id", user)
			if err != nil {
				t.Fatalf("failed to delete user: %v", err)
			}
		}
	}
}

func TestCreateGetStorage(t *testing.T) {
	type testCase struct {
		userIdItemPair []*game.UserItemPair
		expectedRes    []*game.BatchGetStorageRes
	}
	testCases := []testCase{
		{
			userIdItemPair: []*game.UserItemPair{
				{
					UserId: testUserId,
					ItemId: "1",
				}, {
					UserId: testUserId,
					ItemId: "2",
				},
			},
			expectedRes: []*game.BatchGetStorageRes{
				{
					UserId: testUserId,
					ItemData: []*game.StorageData{
						{
							UserId:  testUserId,
							ItemId:  "1",
							Stock:   100,
							IsKnown: true,
						},
						{
							UserId:  testUserId,
							ItemId:  "2",
							Stock:   200,
							IsKnown: true,
						},
					},
				},
			},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		data := func(d []*game.BatchGetStorageRes) []*game.StorageData {
			var res []*game.StorageData
			for _, w := range d {
				for _, x := range w.ItemData {
					res = append(res, x)
				}
			}
			return res
		}(v.expectedRes)
		txErr := dba.Transaction(
			ctx, func(ctx context.Context) error {
				_, err := dba.Exec(
					ctx,
					`INSERT INTO ringo.item_storages (user_id, item_id, stock, is_known) VALUES (:user_id, :item_id, :stock, :is_known)`,
					data,
				)
				if err != nil {
					t.Fatalf("failed to insert storage: %v", err)
				}
				fetchStorage := CreateGetStorage(dba.Query)
				res, err := fetchStorage(ctx, v.userIdItemPair)
				if err != nil {
					t.Fatalf("failed to fetch storage: %v", err)
				}
				if !test.DeepEqual(res, v.expectedRes) {
					t.Errorf("got: %+v, expect: %+v", utils.ToObjArray(res), utils.ToObjArray(v.expectedRes))
				}
				return TestCompleted
			},
		)
		if !errors.Is(txErr, TestCompleted) && txErr != nil {
			t.Fatalf("failed: %v", txErr)
		}
	}
}

func TestCreateGetAllStorage(t *testing.T) {
	type testCase struct {
		userId  core.UserId
		storage []*game.StorageData
	}

	testCases := []testCase{
		{
			userId: testUserId,
			storage: []*game.StorageData{
				{
					UserId: testUserId,
					ItemId: "3",
					Stock:  100,
				},
				{
					UserId: testUserId,
					ItemId: "4",
					Stock:  200,
				},
			},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		_, err := dba.Exec(
			ctx,
			`INSERT INTO ringo.item_storages (user_id, item_id, stock, is_known) VALUES (:user_id, :item_id, :stock, :is_known)`,
			v.storage,
		)
		if err != nil {
			t.Fatalf("failed to insert storage: %v", err)
		}

		fetchAllStorage := CreateGetAllStorage(dba.Query)
		res, err := fetchAllStorage(ctx, v.userId)
		if err != nil {
			t.Fatalf("failed to fetch all storage: %v", err)
		}
		if !test.DeepEqual(res, v.storage) {
			t.Errorf("got: %+v, expect: %+v", res, v.storage)
		}
		_, err = dba.Exec(ctx, `DELETE FROM ringo.item_storages WHERE user_id = :user_id`, v.storage[0])
		if err != nil {
			t.Fatalf("failed to delete storage: %v", err)
		}
	}
}

func TestCreateUpdateItemStorage(t *testing.T) {
	type testCase struct {
		dataBefore []*game.StorageData
		data       []*game.StorageData
	}

	testCases := []testCase{
		{
			data: []*game.StorageData{
				{
					UserId: testUserId,
					ItemId: "1",
					Stock:  100,
				},
				{
					UserId: testUserId,
					ItemId: "2",
					Stock:  200,
				},
			},
			dataBefore: []*game.StorageData{
				{
					UserId: testUserId,
					ItemId: "1",
					Stock:  50,
				},
				{
					UserId: testUserId,
					ItemId: "2",
					Stock:  500,
				},
			},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		_, err := dba.Exec(
			ctx,
			`INSERT INTO ringo.item_storages (user_id, item_id, stock, is_known) VALUES (:user_id, :item_id, :stock, :is_known)`,
			v.dataBefore,
		)
		if err != nil {
			t.Fatalf("failed to insert storage: %v", err)
		}
		err = dba.Transaction(
			ctx, func(ctx context.Context) error {
				updateStorage := CreateUpdateItemStorage(dba.Exec)
				err = updateStorage(ctx, v.data)
				if err != nil {
					t.Fatalf("failed to update storage: %v", err)
				}
				for _, w := range v.data {
					rows, err := dba.Query(
						ctx,
						fmt.Sprintf(
							`SELECT stock FROM ringo.item_storages WHERE user_id = "%s" AND item_id = "%s"`,
							w.UserId,
							w.ItemId,
						),
						nil,
					)
					if err != nil {
						t.Fatalf("failed to fetch storage: %v", err)
					}
					if !rows.Next() {
						t.Fatalf("failed to fetch storage: %v", err)
					}
					var res core.Stock
					err = rows.Scan(&res)
					if err != nil {
						t.Fatalf("failed to fetch storage: %v", err)
					}
					err = rows.Close()
					if err != nil {
						t.Fatalf("failed to close rows: %v", err)
					}
					if res != w.Stock {
						t.Errorf("got: %v, expect: %v", res, w.Stock)
					}
				}
				return TestCompleted
			},
		)
		if !errors.Is(err, TestCompleted) {
			t.Fatalf("failed: %v", err)
		}
	}
}

func TestCreateGetUserSkill(t *testing.T) {
	type testCase struct {
		userId   core.UserId
		skillIds []core.SkillId
		expected game.BatchGetUserSkillRes
	}

	type userSkillData struct {
		UserId   core.UserId   `db:"user_id"`
		SkillId  core.SkillId  `db:"skill_id"`
		SkillExp core.SkillExp `db:"skill_exp"`
	}

	testCases := []testCase{
		{
			userId:   testUserId,
			skillIds: []core.SkillId{"1", "2"},
			expected: game.BatchGetUserSkillRes{
				UserId: testUserId,
				Skills: []*game.UserSkillRes{
					{
						UserId:   testUserId,
						SkillId:  "1",
						SkillExp: 100,
					},
					{
						UserId:   testUserId,
						SkillId:  "2",
						SkillExp: 200,
					},
				},
			},
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		data := func() []*userSkillData {
			var res []*userSkillData
			for _, w := range v.expected.Skills {
				res = append(
					res, &userSkillData{
						UserId:   v.userId,
						SkillId:  w.SkillId,
						SkillExp: w.SkillExp,
					},
				)
			}
			return res
		}()
		txErr := dba.Transaction(
			ctx, func(ctx context.Context) error {
				_, err := dba.Exec(
					ctx,
					`INSERT INTO ringo.user_skills (user_id, skill_id, skill_exp) VALUES (:user_id, :skill_id, :skill_exp)`,
					data,
				)
				if err != nil {
					t.Fatalf("failed to insert user skill: %v", err)
				}
				fetchUserSkill := CreateGetUserSkill(dba.Query)
				res, err := fetchUserSkill(ctx, v.userId, v.skillIds)
				if err != nil {
					t.Fatalf("failed to fetch user skill: %v", err)
				}
				if !test.DeepEqual(res, v.expected) {
					t.Errorf("got: %+v, expect: %+v", res, v.expected)
					if res.UserId != v.expected.UserId {
						t.Errorf("got: %v, expect: %v", res.UserId, v.expected.UserId)
					}
					if len(res.Skills) != len(v.expected.Skills) {
						t.Errorf("got: %d, expect: %d", len(res.Skills), len(v.expected.Skills))
					}
					for i := range res.Skills {
						r := res.Skills[i]
						e := v.expected.Skills[i]
						if r.SkillId != e.SkillId {
							t.Errorf("got: %v, expect: %v", r.SkillId, e.SkillId)
						}
						if r.SkillExp != e.SkillExp {
							t.Errorf("got: %v, expect: %v", r.SkillExp, e.SkillExp)
						}
					}
				}
				return TestCompleted
			},
		)
		if !errors.Is(txErr, TestCompleted) && txErr != nil {
			t.Fatalf("failed: %v", txErr)
		}
	}
}

func TestCreateUpdateUserSkill(t *testing.T) {
	type testCase struct {
		userId   core.UserId
		skillId  core.SkillId
		skillExp core.SkillExp
		expAfter core.SkillExp
	}

	testCases := []testCase{
		{
			userId:   testUserId,
			skillId:  "1",
			skillExp: 100,
			expAfter: 200,
		},
	}

	for _, v := range testCases {
		ctx := test.MockCreateContext()
		data := func() *game.UserSkillRes {
			return &game.UserSkillRes{
				UserId:   v.userId,
				SkillId:  v.skillId,
				SkillExp: v.skillExp,
			}
		}()
		txErr := dba.Transaction(
			ctx, func(tx context.Context) error {
				_, err := dba.Exec(
					ctx,
					`INSERT INTO ringo.user_skills (user_id, skill_id, skill_exp) VALUES (:user_id, :skill_id, :skill_exp)`,
					data,
				)
				if err != nil {
					t.Fatalf("failed to insert user skill: %v", err)
					return err
				}
				updateUserSkill := CreateUpdateUserSkill(dba.Exec)
				err = updateUserSkill(
					ctx, game.SkillGrowthPost{
						UserId: v.userId,
						SkillGrowth: []*game.SkillGrowthPostRow{
							{
								UserId:   v.userId,
								SkillId:  v.skillId,
								SkillExp: v.expAfter,
							},
						},
					},
				)
				if err != nil {
					t.Fatalf("failed to update user skill: %v", err)
					return err
				}
				rows, err := dba.Query(
					ctx,
					fmt.Sprintf(
						`SELECT skill_exp FROM ringo.user_skills WHERE user_id = "%s" AND skill_id = "%s"`,
						v.userId,
						v.skillId,
					),
					nil,
				)
				if err != nil {
					t.Fatalf("failed to fetch user skill: %v", err)
					return err
				}
				if !rows.Next() {
					t.Fatalf("failed to fetch user skill: %v", err)
					return err
				}
				var res core.SkillExp
				err = rows.Scan(&res)
				if err != nil {
					t.Fatalf("failed to fetch user skill: %v", err)
					return err
				}
				if res != v.expAfter {
					t.Errorf("got: %v, expect: %v", res, v.expAfter)
				}
				return TestCompleted
			},
		)
		if !errors.Is(txErr, TestCompleted) {
			t.Errorf("failed: %v", txErr)
		}
	}
}
