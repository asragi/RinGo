package infrastructure

import (
	"testing"

	"github.com/asragi/RinGo/stage"
)

type loader struct{}

func (loader *loader) Load() (ExploreMasterData, error) {
	return ExploreMasterData{
		stage.ExploreId("1"): stage.GetExploreMasterRes{
			ExploreId: stage.ExploreId("1"),
		},
	}, nil
}

func TestExploreMasterRepo(t *testing.T) {
	repo, err := CreateInMemoryExploreMasterRepo(&loader{})
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	res, _ := repo.BatchGet([]stage.ExploreId{stage.ExploreId("1")})

	if len(res) != 1 {
		t.Errorf("res length: %d", len(res))
	}
}
