package auth

import "github.com/asragi/RinGo/core"

type InsertNewUser func(core.UserId, hashedPassword) error
