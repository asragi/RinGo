package auth

type rowPassword string
type createRowPasswordFunc func() rowPassword
type hashedPassword string
type secretHashKey string
type createHashedPasswordFunc func(rowPassword) (hashedPassword, error)

type EncryptFunc func(string) (string, error)

func createHashedPassword(encrypt EncryptFunc) createHashedPasswordFunc {
	return func(password rowPassword) (hashedPassword, error) {
		passwordString, err := encrypt(string(password))
		return hashedPassword(passwordString), err
	}
}
