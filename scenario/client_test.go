package scenario

import (
	"fmt"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RinGo/core/game/explore"
	"github.com/asragi/RinGo/utils"
	"github.com/asragi/RingoSuPBGo/gateway"
	"google.golang.org/grpc"
)

type client struct {
	connectFunc     ConnectFunc
	userId          core.UserId
	password        auth.RowPassword
	token           auth.AccessToken
	stageInfo       *utils.Set[*gateway.StageInformation]
	itemList        *utils.Set[*gateway.GetItemListResponseRow]
	itemDetailCache *gateway.GetItemDetailResponse
	shelves         *utils.Set[*gateway.Shelf]
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

func (c *client) storeStageInfo(info []*gateway.StageInformation) {
	c.stageInfo = utils.NewSet(info)
}

func (c *client) pickStageAction() (explore.StageId, game.ActionId, error) {
	type stageExplorePair struct {
		stageId   explore.StageId
		exploreId game.ActionId
	}
	var allError error
	var possibleActions []*stageExplorePair
	c.stageInfo.Foreach(
		func(_ int, info *gateway.StageInformation) {
			if allError != nil {
				return
			}
			for _, e := range info.UserExplore {
				if !e.IsPossible {
					continue
				}
				actionId, err := game.NewActionId(e.ExploreId)
				allError = err
				if err != nil {
					return
				}
				stageId, err := explore.CreateStageId(info.StageId)
				allError = err
				if err != nil {
					return
				}
				possibleActions = append(
					possibleActions, &stageExplorePair{
						stageId:   stageId,
						exploreId: actionId,
					},
				)
			}
		},
	)
	if allError != nil {
		return *new(explore.StageId), *new(game.ActionId), allError
	}
	targetAction := possibleActions[0]
	return targetAction.stageId, targetAction.exploreId, nil
}

func (c *client) storeItemList(itemList []*gateway.GetItemListResponseRow) {
	c.itemList = utils.NewSet(itemList)
}

func (c *client) selectItem() (core.ItemId, error) {
	var emptyItemId core.ItemId
	if c.itemList.Length() == 0 {
		return emptyItemId, fmt.Errorf("item list is empty")
	}
	item := c.itemList.Get(0)
	itemId, err := core.NewItemId(item.ItemId)
	if err != nil {
		return emptyItemId, err
	}
	return itemId, nil
}

func (c *client) storeItemDetail(response *gateway.GetItemDetailResponse) {
	c.itemDetailCache = response
}

func (c *client) getItemAction() (core.ItemId, game.ActionId, error) {
	handleError := func(err error) (core.ItemId, game.ActionId, error) {
		return "", "", fmt.Errorf("get item action: %w", err)
	}
	if c.itemDetailCache == nil {
		return handleError(fmt.Errorf("item detail is not stored"))
	}
	itemId, err := core.NewItemId(c.itemDetailCache.ItemId)
	if err != nil {
		return handleError(err)
	}
	actionSet := utils.NewSet(c.itemDetailCache.UserExplore)
	if actionSet.Length() == 0 {
		return handleError(fmt.Errorf("no action"))
	}
	possibleAction := actionSet.Find(
		func(e *gateway.UserExplore) bool {
			return e.IsPossible
		},
	)
	if possibleAction == nil {
		return handleError(fmt.Errorf("no possible action"))
	}
	actionId, err := game.NewActionId(possibleAction.ExploreId)
	if err != nil {
		return handleError(err)
	}
	return itemId, actionId, nil
}

func (c *client) storeShelves(shelves []*gateway.Shelf) {
	c.shelves = utils.NewSet(shelves)
}

func (c *client) selectShelf() *gateway.Shelf {
	return c.shelves.Get(0)
}
