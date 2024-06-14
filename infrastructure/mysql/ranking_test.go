package mysql

import (
	"context"
	"errors"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game/shelf/ranking"
	"github.com/asragi/RinGo/debug"
	"github.com/asragi/RinGo/test"
	"testing"
)

func TestCreateInsertRankPeriod(t *testing.T) {
	period := ranking.RankPeriod(2)
	ctx := test.MockCreateContext()
	insertRank := CreateInsertRankPeriod(dba.Exec)
	txErr := dba.Transaction(
		ctx,
		func(ctx context.Context) error {
			err := insertRank(ctx, period)
			if err != nil {
				return err
			}
			return TestCompleted
		},
	)
	if !errors.Is(txErr, TestCompleted) {
		t.Errorf("transaction error: %v", txErr)
	}
}

func TestCreateInsertWin(t *testing.T) {
	secondUser := core.UserId("second_user")
	req := []*ranking.InsertWinReq{
		{
			UserId: testUserId,
			Rank:   1,
			Period: 1,
		},
		{
			UserId: secondUser,
			Rank:   2,
			Period: 1,
		},
	}
	insertWin := CreateInsertWin(dba.Exec)
	ctx := test.MockCreateContext()
	txErr := dba.Transaction(
		ctx,
		func(ctx context.Context) error {
			err := addTestUser(
				func(u *userTest) {
					u.UserId = secondUser
				},
			)
			if err != nil {
				return err
			}
			err = debug.CreateAddInitialPeriod(dba.Exec)(ctx)
			if err != nil {
				return err
			}
			err = insertWin(ctx, req)
			if err != nil {
				return err
			}
			return TestCompleted
		},
	)
	if !errors.Is(txErr, TestCompleted) {
		t.Errorf("transaction error: %v", txErr)
	}
}
