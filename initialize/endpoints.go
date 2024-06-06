package initialize

import (
	"github.com/asragi/RinGo/endpoint"
	"github.com/google/wire"
)

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
