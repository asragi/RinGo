package auth

import "github.com/asragi/RinGo/core"

type InsertNewUser func(*core.UserId, *core.UserName, *HashedPassword) error
type FetchHashedPassword func(*core.UserId) (*HashedPassword, error)
