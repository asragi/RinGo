package scenario

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type adminClient struct {
	connector ConnectFunc
	token     auth.AccessToken
}

func newAdmin(address string) *adminClient {
	return &adminClient{
		connector: Connect(address),
	}
}

func (a *adminClient) saveToken(token auth.AccessToken) {
	a.token = token
}

func (a *adminClient) useToken() auth.AccessToken {
	return a.token
}

func (a *adminClient) getDebugTimeClient() (gateway.DebugTimeClient, closeConnectionType, error) {
	conn, err := a.connector()
	if err != nil {
		return nil, nil, err
	}
	closeConnWrapper := func() {
		closeConnection(conn)
	}
	return gateway.NewDebugTimeClient(conn), closeConnWrapper, nil
}

func (a *adminClient) getChangePeriodClient() (gateway.ChangePeriodClient, closeConnectionType, error) {
	conn, err := a.connector()
	if err != nil {
		return nil, nil, err
	}
	closeConnWrapper := func() {
		closeConnection(conn)
	}
	return gateway.NewChangePeriodClient(conn), closeConnWrapper, nil
}

func (a *adminClient) useLoginData() (core.UserId, auth.RowPassword) {
	return "admin", "admin"
}
