package initialize

import (
	"github.com/asragi/RinGo/endpoint"
	"github.com/google/wire"
)

type Endpoints struct {
	SignUp               endpoint.RegisterEndpointFunc
	Login                endpoint.LoginEndpoint
	UpdateUserName       endpoint.UpdateUserNameEndpoint
	UpdateShopName       endpoint.UpdateShopNameEndpoint
	GetResource          endpoint.GetResourceEndpoint
	GetItemList          endpoint.GetItemListEndpoint
	GetItemDetail        endpoint.GetItemDetailEndpointFunc
	GetItemActionDetail  endpoint.GetItemActionDetailEndpoint
	GetMyShelves         endpoint.GetMyShelvesEndpointFunc
	GetRankingUserList   endpoint.GetRankingUserListEndpoint
	GetStageList         endpoint.GetStageListEndpointFunc
	GetStageActionDetail endpoint.GetStageActionEndpointFunc
	PostAction           endpoint.PostActionEndpointFunc
	UpdateShelfContent   endpoint.UpdateShelfContentEndpointFunc
	UpdateShelfSize      endpoint.UpdateShelfSizeEndpoint
}

var endpointsSet = wire.NewSet(
	endpoint.CreateRegisterEndpoint,
	endpoint.CreateLoginEndpoint,
	endpoint.CreateUpdateUserNameEndpoint,
	endpoint.CreateUpdateShopNameEndpoint,
	endpoint.CreateGetResourceEndpoint,
	endpoint.CreateGetItemService,
	endpoint.CreateGetItemDetail,
	endpoint.CreateGetItemActionDetailEndpoint,
	endpoint.CreateGetMyShelvesEndpoint,
	endpoint.CreateGetRankingUserList,
	endpoint.CreateGetStageList,
	endpoint.CreateGetStageActionDetail,
	endpoint.CreatePostAction,
	endpoint.CreateUpdateShelfContentEndpoint,
	endpoint.CreateUpdateShelfSizeEndpoint,
)
