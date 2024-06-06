package mysql

import (
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/test"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

var dba *database.DBAccessor
var testUserId = core.UserId("the-one-test-user")

func TestMain(m *testing.M) {
	db, purge, err := test.CreateTestDB("ringo-mysql-unittest-image", "../../test/db_for_test/Dockerfile")
	if err != nil {
		log.Fatalf("Could not create test DB: %s", err)
		return
	}
	defer func() {
		if err = purge(); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	dba = database.NewDBAccessor(db, db)

	err = addTestUser(func(u *userTest) { u.UserId = testUserId })
	if err != nil {
		log.Fatalf("Could not add test user: %s", err)
		return
	}

	m.Run()
}

var insertTestUserQuery = "INSERT INTO ringo.users (user_id, name, shop_name, max_stamina, stamina_recover_time, fund, popularity, hashed_password) VALUES (:user_id, :name, :shop_name, :max_stamina, :stamina_recover_time, :fund, :popularity, :hashed_password)"

func addTestUser(options ...ApplyUserTestOption) error {
	user := createTestUser(options...)
	ctx := test.MockCreateContext()
	_, err := dba.Exec(
		ctx,
		insertTestUserQuery,
		user,
	)
	if err != nil {
		return err
	}
	return nil
}
