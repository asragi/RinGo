package endpoint

import (
	"context"
	"fmt"
	"github.com/asragi/RinGo/auth"

	"github.com/asragi/RinGo/application"
	"github.com/asragi/RinGo/stage"
	"github.com/asragi/RingoSuPBGo/gateway"
)

type CreatePostActionEndpoint func(application.CreatePostActionRes, auth.ValidateTokenFunc) postActionEndpoint
type postActionEndpoint func(context.Context, *gateway.PostActionRequest) (*gateway.PostActionResponse, error)

func CreatePostAction(
	postAction application.CreatePostActionRes,
	validateToken auth.ValidateTokenFunc,
) postActionEndpoint {
	post := func(ctx context.Context, req *gateway.PostActionRequest) (*gateway.PostActionResponse, error) {
		handleError := func(err error) (*gateway.PostActionResponse, error) {
			return &gateway.PostActionResponse{
				Error: &gateway.Error{
					ErrorOccurred:  true,
					DisplayMessage: err.Error(),
				},
			}, fmt.Errorf("error on post action: %w", err)
		}
		exploreId := stage.ExploreId(req.ExploreId)
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return handleError(err)
		}
		userId := tokenInfo.UserId
		execCount := int(req.ExecCount)
		res, err := postAction(ctx, userId, token, exploreId, execCount)
		if err != nil {
			return handleError(err)
		}

		earnedItem := func() []*gateway.EarnedItems {
			result := make([]*gateway.EarnedItems, len(res.EarnedItems))
			for i, v := range res.EarnedItems {
				result[i] = &gateway.EarnedItems{
					ItemId: string(v.ItemId),
					Count:  int32(v.Count),
				}
			}
			return result
		}()
		consumedItem := func() []*gateway.ConsumedItems {
			result := make([]*gateway.ConsumedItems, len(res.ConsumedItems))
			for i, v := range res.ConsumedItems {
				result[i] = &gateway.ConsumedItems{
					ItemId: string(v.ItemId),
					Count:  int32(v.Count),
				}
			}
			return result
		}()
		skillGrowth := func() []*gateway.SkillGrowthResult {
			result := make([]*gateway.SkillGrowthResult, len(res.SkillGrowthInformation))
			for i, v := range res.SkillGrowthInformation {
				result[i] = &gateway.SkillGrowthResult{
					DisplayName: string(v.DisplayName),
					BeforeExp:   int32(v.GrowthResult.BeforeExp),
					BeforeLv:    int32(v.GrowthResult.BeforeLv),
					SkillId:     string(v.GrowthResult.SkillId),
					AfterExp:    int32(v.GrowthResult.AfterExp),
					AfterLv:     int32(v.GrowthResult.AfterLv),
				}
			}
			return result
		}()
		return &gateway.PostActionResponse{
			Error: &gateway.Error{
				ErrorOccurred:  false,
				DisplayMessage: "",
			},
			EarnedItems:       earnedItem,
			ConsumedItems:     consumedItem,
			SkillGrowthResult: skillGrowth,
		}, nil
	}

	return post
}
