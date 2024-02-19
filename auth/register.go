package auth

import (
	"fmt"
	"github.com/asragi/RinGo/core"
)

type registerResult struct {
	UserId   core.UserId
	Password RowPassword
}

type generateIdStringFunc func() string

type createUserIdFunc func() (core.UserId, error)

func CreateUserId(
	challengeNum int,
	checkUser core.CheckDoesUserExist,
	generate generateIdStringFunc,
) createUserIdFunc {
	f := func() (core.UserId, error) {
		var err error
		for i := 0; i < challengeNum; i++ {
			userId := core.UserId(generate())
			err = checkUser(userId)
			if err == nil {
				return userId, nil
			}
		}
		return "", fmt.Errorf("creating user id was failed: %w", err)
	}
	return f
}

type decideInitialName func() core.UserName
type RegisterUserFunc func() (registerResult, error)

func RegisterUser(
	generateUserId createUserIdFunc,
	generateRowPassword createRowPasswordFunc,
	createHashedPassword createHashedPasswordFunc,
	insertNewUser InsertNewUser,
	decideName decideInitialName,
) RegisterUserFunc {
	f := func() (registerResult, error) {
		handleError := func(err error) (registerResult, error) {
			return registerResult{}, fmt.Errorf("register user: %w", err)
		}
		userId, err := generateUserId()
		if err != nil {
			return handleError(err)
		}
		rowPass := generateRowPassword()
		hashedPass, err := createHashedPassword(rowPass)
		if err != nil {
			return handleError(err)
		}
		initialName := decideName()
		err = insertNewUser(&userId, &initialName, &hashedPass)
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
