package core

import (
	"testing"
)

type MockItemRepo struct{}

var mockItemId ItemId = "Item-Ringo"

func (m *MockItemRepo) Get(
	userId UserId,
	itemId ItemId,
	token AccessToken) (GetItemDetailRes, error) {
	return GetItemDetailRes{
		ItemID: mockItemId,
	}, nil
}

func TestCreateItemService(t *testing.T) {
	itemRepo := MockItemRepo{}
	itemService := CreateItemService(&itemRepo)
	getUserItemDetail := itemService.GetUserItemDetail
	req := GetUserItemDetailReq{}
	res := getUserItemDetail(req)

	if res.ItemId != ItemId(mockItemId) {
		t.Errorf("want %s, actual %s", mockItemId, res.ItemId)
	}
}
