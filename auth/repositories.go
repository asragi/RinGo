package auth

import "github.com/asragi/RinGo/core"

type InsertNewUser func(core.UserId, hashedPassword) error
type FetchHashedPassword func(core.UserId) (hashedPassword, error)

type ValidateTokenRepoFunc func(core.UserId, AccessToken) error
