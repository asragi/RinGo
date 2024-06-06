package endpoint

type Endpoints struct {
	SignUp               RegisterEndpointFunc
	Login                LoginEndpoint
	UpdateUserName       UpdateUserNameEndpoint
	UpdateShopName       UpdateShopNameEndpoint
	GetResource          GetResourceEndpoint
	GetItemList          GetItemListEndpoint
	GetItemDetail        GetItemDetailEndpointFunc
	GetItemActionDetail  GetItemActionDetailEndpoint
	GetMyShelves         GetMyShelvesEndpointFunc
	GetRankingUserList   GetRankingUserListEndpoint
	GetStageList         GetStageListEndpointFunc
	GetStageActionDetail GetStageActionEndpointFunc
	PostAction           PostActionEndpointFunc
	UpdateShelfContent   UpdateShelfContentEndpointFunc
	UpdateShelfSize      UpdateShelfSizeEndpoint
}
