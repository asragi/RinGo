//go:build wireinject
// +build wireinject

package initialize

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/database"
	"github.com/asragi/RinGo/endpoint"
	"github.com/google/wire"
)

func CreateEndpoints(
	key auth.SecretHashKey,
	constants *Constants,
	exec database.ExecFunc,
	query database.QueryFunc,
) *endpoint.Endpoints {
	wire.Build(
		wire.FieldsOf(
			new(*Constants),
			"UserIdChallengeNum",
			"InitialFund",
			"InitialMaxStamina",
			"InitialPopularity",
		),
		commonSet,
		infrastructures,
		services,
		endpointsSet,
		wire.Struct(new(endpoint.Endpoints), "*"),
	)
	return nil
}
