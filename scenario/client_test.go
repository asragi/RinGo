package scenario

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RingoSuPBGo/gateway"
	"google.golang.org/grpc"
)

type client struct {
	connectFunc ConnectFunc
	userId      core.UserId
	password    auth.RowPassword
	token       auth.AccessToken
}

type closeConnectionType func()
type connectAgent interface {
	// deprecated: use getClient
	connect() (*grpc.ClientConn, error)
	getClient() (gateway.RingoClient, closeConnectionType, error)
}

type useToken interface {
	useToken() auth.AccessToken
}

func newClient(address string) *client {
	return &client{
		connectFunc: Connect(address),
	}
}

func (c *client) connect() (*grpc.ClientConn, error) {
	return c.connectFunc()
}

func (c *client) getClient() (gateway.RingoClient, closeConnectionType, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect: %w", err)
	}
	closeConnWrapper := func() {
		closeConnection(conn)
	}
	return gateway.NewRingoClient(conn), closeConnWrapper, nil
}

func (c *client) saveUserData(userId core.UserId, password auth.RowPassword) {
	c.userId = userId
	c.password = password
}

func (c *client) saveToken(token auth.AccessToken) {
	c.token = token
}

func (c *client) useLoginData() (core.UserId, auth.RowPassword) {
	return c.userId, c.password
}

func (c *client) useToken() auth.AccessToken {
	return c.token
}
