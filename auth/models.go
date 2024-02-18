package auth

type RowPassword string
type createRowPasswordFunc func() RowPassword
type hashedPassword string
type SecretHashKey string
type createHashedPasswordFunc func(RowPassword) (hashedPassword, error)

type EncryptFunc func(string) (string, error)

func createHashedPassword(encrypt EncryptFunc) createHashedPasswordFunc {
	return func(password RowPassword) (hashedPassword, error) {
		passwordString, err := encrypt(string(password))
		return hashedPassword(passwordString), err
	}
}
