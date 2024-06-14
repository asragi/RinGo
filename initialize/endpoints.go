package initialize

import (
	"github.com/asragi/RinGo/endpoint"
	"github.com/google/wire"
)

var devEndpointsSet = wire.NewSet(
	endpoint.CreateAdminLoginEndpoint,
	endpoint.CreateAutoInsertReservationEndpoint,
	endpoint.CreateChangePeriod,
	endpoint.CreateChangeTimeEndpoint,
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
	devEndpointsSet,
)
