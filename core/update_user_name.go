package core

import "context"

type UpdateUserNameServiceFunc func(context.Context, UserId, Name) error

func CreateUpdateUserNameServiceFunc(updateUserName UpdateUserNameFunc) UpdateUserNameServiceFunc {
	return func(ctx context.Context, userId UserId, userName Name) error {
		return updateUserName(ctx, userId, userName)
	}
}
