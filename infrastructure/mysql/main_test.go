package mysql

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/test"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"log"
	"testing"
)

var dba *database.DBAccessor
var testUserId = core.UserId("the-one-test-user")

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.BuildAndRun("mysql-test-image", "./db_for_test/Dockerfile", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	// err = resource.Expire(20)
	if err != nil {
		log.Fatalf("Could not set expiration time: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(
		func() error {
			var err error
			db, err := sqlx.Open(
				"mysql",
				fmt.Sprintf("root:ringo@(localhost:%s)/mysql?parseTime=true", resource.GetPort("3306/tcp")),
			)
			if err != nil {
				return err
			}
			dba = database.NewDBAccessor(db, db)
			return db.Ping()
		},
	); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	err = addTestUser(func(u *userTest) { u.UserId = testUserId })
	if err != nil {
		log.Fatalf("Could not add test user: %s", err)
	}

	m.Run()
}

func TestSomething(t *testing.T) {
	fmt.Printf("TEST IS HERE!")
}

var insertTestUserQuery = "INSERT INTO ringo.users (user_id, name, max_stamina, stamina_recover_time, fund, popularity) VALUES (:user_id, :name, :max_stamina, :stamina_recover_time, :fund, :popularity)"

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
