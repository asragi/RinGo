package core

// CheckDoesUserExist returns error when user_id is already used
type CheckDoesUserExist func(UserId) error
type ValidateTokenRepoFunc func(UserId, AccessToken) error
