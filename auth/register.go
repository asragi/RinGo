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

type generateIdStringFunc func() string

type createUserIdFunc func(context.Context) (core.UserId, error)

func CreateUserId(
	challengeNum int,
	checkUser core.CheckDoesUserExist,
	generate generateIdStringFunc,
) createUserIdFunc {
	f := func(ctx context.Context) (core.UserId, error) {
		var err error
		for i := 0; i < challengeNum; i++ {
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

type decideInitialName func() core.Name
type RegisterUserFunc func(context.Context) (registerResult, error)

func RegisterUser(
	generateUserId createUserIdFunc,
	generateRowPassword createRowPasswordFunc,
	createHashedPassword createHashedPasswordFunc,
	insertNewUser InsertNewUser,
	decideName decideInitialName,
	decideShopName decideInitialName,
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
