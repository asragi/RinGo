package admin

import (
	"context"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
)

type CreateCommonLoginFunc func(
	fetchHashedPassword auth.FetchHashedPassword,
	comparePassword auth.CompareHashedPassword,
	createToken auth.CreateTokenFunc,
) auth.LoginFunc

type LoginFunc func(context.Context, core.UserId, auth.RowPassword) (auth.AccessToken, error)

func CreateLogin(
	fetchHashedPassword FetchHashedPassword,
	comparePassword auth.CompareHashedPassword,
	createToken auth.CreateTokenFunc,
	createLogin CreateCommonLoginFunc,
) LoginFunc {
	fetchPassword := auth.FetchHashedPassword(fetchHashedPassword)
	f := createLogin(fetchPassword, comparePassword, createToken)
	return LoginFunc(f)
}
