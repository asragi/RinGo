package auth

type RowPassword string

func NewRowPassword(password string) RowPassword {
	return RowPassword(password)
}

func (p RowPassword) String() string {
	return string(p)
}

type RowPasswordGenerator func() string

func GenerateRowPassword(gen RowPasswordGenerator) CreateRowPasswordFunc {
	return func() RowPassword { return NewRowPassword(gen()) }
}

type CreateRowPasswordFunc func() RowPassword
type HashedPassword string
type SecretHashKey string
type CreateHashedPasswordFunc func(RowPassword) (HashedPassword, error)

type EncryptFunc func(string) (string, error)

func CreateHashedPassword(encrypt EncryptFunc) CreateHashedPasswordFunc {
	return func(password RowPassword) (HashedPassword, error) {
		passwordString, err := encrypt(string(password))
		return HashedPassword(passwordString), err
	}
}
