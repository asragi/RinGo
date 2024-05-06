package core

type Service struct {
	UpdateUserName UpdateUserNameServiceFunc
	UpdateShopName UpdateShopNameServiceFunc
}

func NewService(updateUserName UpdateUserNameFunc, updateShopName UpdateShopNameFunc) *Service {
	return &Service{
		UpdateUserName: CreateUpdateUserNameServiceFunc(updateUserName),
		UpdateShopName: CreateUpdateShopNameServiceFunc(updateShopName),
	}
}
