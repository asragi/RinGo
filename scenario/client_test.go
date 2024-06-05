package scenario

import (
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"google.golang.org/grpc"
)

type client struct {
	connectFunc ConnectFunc
	userId      core.UserId
	password    auth.RowPassword
	token       auth.AccessToken
}

func newClient(address string) *client {
	return &client{
		connectFunc: Connect(address),
	}
}

func (c *client) connect() (*grpc.ClientConn, error) {
	return c.connect()
}

func (c *client) saveUserData(userId core.UserId, password auth.RowPassword) {
	c.userId = userId
	c.password = password
}
