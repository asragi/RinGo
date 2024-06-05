package auth

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/core"
)

type registerResult struct {
	UserId   core.UserId
	Password RowPassword
}

type createUserIdFunc func(context.Context) (core.UserId, error)

type CreateUserIdChallengeNum int

func CreateUserId(
	challengeNum CreateUserIdChallengeNum,
	checkUser core.CheckDoesUserExist,
	generate core.GenerateUUIDFunc,
) createUserIdFunc {
	f := func(ctx context.Context) (core.UserId, error) {
		var err error
		for i := 0; i < int(challengeNum); i++ {
			userId := core.UserId(generate())
			err = checkUser(ctx, userId)
			if err == nil {
				return userId, nil
			}
		}
		return "", fmt.Errorf("creating user id was failed: %w", err)
	}
	return f
}

type RegisterUserFunc func(context.Context) (registerResult, error)

func RegisterUser(
	generateUserId createUserIdFunc,
	generateRowPassword CreateRowPasswordFunc,
	createHashedPassword createHashedPasswordFunc,
	insertNewUser InsertNewUser,
	decideName core.DecideInitialName,
	decideShopName core.DecideInitialShopName,
) RegisterUserFunc {
	f := func(ctx context.Context) (registerResult, error) {
		handleError := func(err error) (registerResult, error) {
			return registerResult{}, fmt.Errorf("register user: %w", err)
		}
		userId, err := generateUserId(ctx)
		if err != nil {
			return handleError(err)
		}
		rowPass := generateRowPassword()
		hashedPass, err := createHashedPassword(rowPass)
		if err != nil {
			return handleError(err)
		}
		initialName := decideName()
		initialShopName := decideShopName()
		err = insertNewUser(ctx, userId, initialName, initialShopName, hashedPass)
		if err != nil {
			return handleError(err)
		}
		return registerResult{
			UserId:   userId,
			Password: rowPass,
		}, nil
	}
	return f
}
