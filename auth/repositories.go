package auth

import "github.com/asragi/RinGo/core"

type InsertNewUser func(*core.UserId, *HashedPassword) error
type FetchHashedPassword func(*core.UserId) (*HashedPassword, error)
