package auth

type RowPassword string
type createRowPasswordFunc func() RowPassword
type HashedPassword string
type SecretHashKey string
type createHashedPasswordFunc func(RowPassword) (HashedPassword, error)

type EncryptFunc func(string) (string, error)

func createHashedPassword(encrypt EncryptFunc) createHashedPasswordFunc {
	return func(password RowPassword) (HashedPassword, error) {
		passwordString, err := encrypt(string(password))
		return HashedPassword(passwordString), err
	}
}
