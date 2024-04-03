package infrastructure

import (
	"errors"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/location"
	"github.com/asragi/RinGo/test"
	"testing"
	"time"
)

type userTest struct {
	UserId             core.UserId         `db:"user_id"`
	Name               core.UserName       `db:"name"`
	MaxStamina         core.MaxStamina     `db:"max_stamina"`
	Fund               core.Fund           `db:"fund"`
	StaminaRecoverTime time.Time           `db:"stamina_recover_time"`
	HashedPassword     auth.HashedPassword `db:"hashed_password"`
}

type ApplyUserTestOption func(*userTest)

func createTestUser(options ...ApplyUserTestOption) *userTest {
	user := userTest{
		UserId:             "test-user",
		Name:               "test-name",
		MaxStamina:         6000,
		Fund:               100000,
		StaminaRecoverTime: test.MockTime(),
		HashedPassword:     "test-password",
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
		"INSERT INTO ringo.users (user_id, name, max_stamina, stamina_recover_time, fund) VALUES (:user_id, :name, :max_stamina, :stamina_recover_time, :fund)",
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
		{userId: "test-user", hashedPassword: "test-password"},
	}

	for _, v := range testCases {
		user := createTestUser(func(user *userTest) { user.HashedPassword = v.hashedPassword; user.UserId = v.userId })
		ctx := test.MockCreateContext()
		_, err := dba.Exec(
			ctx,
			"INSERT INTO ringo.users (user_id, name, max_stamina, stamina_recover_time, fund, hashed_password) VALUES (:user_id, :name, :max_stamina, :stamina_recover_time, :fund, :hashed_password)",
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

func TestCreateGetResourceMySQL(t *testing.T) {
	type testCase struct {
		UserId             core.UserId     `db:"user_id"`
		MaxStamina         core.MaxStamina `db:"max_stamina"`
		StaminaRecoverTime time.Time       `db:"stamina_recover_time"`
		Fund               core.Fund       `db:"fund"`
	}

	testCases := []testCase{
		{
			UserId:             "test-user",
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
			"INSERT INTO ringo.users (user_id, name, max_stamina, stamina_recover_time, fund) VALUES (:user_id, :name, :max_stamina, :stamina_recover_time, :fund)",
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
			UserId:             "test-user",
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
			"INSERT INTO ringo.users (user_id, name, max_stamina, stamina_recover_time, fund) VALUES (:user_id, :name, :max_stamina, :stamina_recover_time, :fund)",
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
