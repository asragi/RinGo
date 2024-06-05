package core

type Services struct {
	UpdateUserName UpdateUserNameServiceFunc
	UpdateShopName UpdateShopNameServiceFunc
}

func NewService(updateUserName UpdateUserNameFunc, updateShopName UpdateShopNameFunc) *Services {
	return &Services{
		UpdateUserName: CreateUpdateUserNameService(updateUserName),
		UpdateShopName: CreateUpdateShopNameService(updateShopName),
	}
}
