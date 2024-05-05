package core

import "context"

type UpdateUserNameServiceFunc func(context.Context, UserId, UserName) error

func CreateUpdateUserNameServiceFunc(updateUserName UpdateUserNameFunc) UpdateUserNameServiceFunc {
	return func(ctx context.Context, userId UserId, userName UserName) error {
		return updateUserName(ctx, userId, userName)
	}
}
