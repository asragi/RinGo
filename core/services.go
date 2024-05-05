package core

type Service struct {
	UpdateUserName UpdateUserNameServiceFunc
}

func NewService(updateUserName UpdateUserNameFunc) *Service {
	return &Service{
		UpdateUserName: CreateUpdateUserNameServiceFunc(updateUserName),
	}
}
