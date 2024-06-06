package scenario

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/initialize"
	"github.com/asragi/RinGo/server"
	"github.com/asragi/RinGo/test"
	"log"
	"testing"
	"time"
)

var port = 4445

func TestMain(m *testing.M) {
	secretKey := auth.SecretHashKey("secret")
	constants := &initialize.Constants{
		InitialFund:        core.Fund(100000),
		InitialMaxStamina:  core.MaxStamina(6000),
		InitialPopularity:  shelf.ShopPopularity(0),
		UserIdChallengeNum: 3,
	}
	db, purge, err := test.CreateTestDB("ringo-mysql-scenario-test-image", "../test/db_for_test/Dockerfile")
	if err != nil {
		log.Fatalf("Could not create test DB: %s", err)
		return
	}
	defer func() {
		if err = purge(); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()
	dba := database.NewDBAccessor(db, db)
	endpoints := initialize.CreateEndpoints(secretKey, constants, dba.Exec, dba.Query)
	serve, stopDB, err := server.SetUpServer(port, endpoints)
	if err != nil {
		log.Fatalf("Could not set up server: %s", err)
		return
	}
	defer stopDB()
	go func() {
		err = serve()
		if err != nil {
			log.Printf("Http Server Error: %v", err)
			return
		}
	}()
	time.Sleep(1 * time.Second)
	m.Run()
	// teardown
}

func TestE2E(t *testing.T) {
	c := newClient(fmt.Sprintf("localhost:%d", port))
	ctx := context.Background()
	err := signUp(ctx, c)
	if err != nil {
		t.Errorf("sign up: %v", err)
	}
	err = login(ctx, c)
	if err != nil {
		t.Errorf("login: %v", err)
	}
	err = updateUserName(ctx, c)
	if err != nil {
		t.Errorf("update user name: %v", err)
	}
	err = updateShopName(ctx, c)
	if err != nil {
		t.Errorf("update shop name: %+v", err)
	}
	err = getResource(ctx, c)
	if err != nil {
		t.Errorf("get resource: %v", err)
	}
	err = getMyShelves(ctx, c)
	if err != nil {
		t.Errorf("get my shelves: %v", err)
	}
	err = getItemList(ctx, c)
	if err != nil {
		t.Errorf("get item list: %v", err)
	}
	err = getStageList(ctx, c)
	if err != nil {
		t.Errorf("get stage list: %v", err)
	}
}
