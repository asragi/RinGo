package endpoint

import (
	"context"
	"github.com/asragi/RinGo/auth"
	"github.com/asragi/RinGo/core/game"
	"github.com/asragi/RingoSuPBGo/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostActionEndpointFunc func(context.Context, *gateway.PostActionRequest) (*gateway.PostActionResponse, error)

func CreatePostAction(
	postAction game.PostActionFunc,
	validateToken auth.ValidateTokenFunc,
) PostActionEndpointFunc {
	post := func(ctx context.Context, req *gateway.PostActionRequest) (*gateway.PostActionResponse, error) {
		exploreId := game.ExploreId(req.ExploreId)
		token := auth.AccessToken(req.Token)
		tokenInfo, err := validateToken(&token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "error on post action: %+v", err)
		}
		userId := tokenInfo.UserId
		execCount := int(req.ExecCount)
		res, err := postAction(ctx, userId, execCount, exploreId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error on post action: %+v", err)
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
			EarnedItems:       earnedItem,
			ConsumedItems:     consumedItem,
			SkillGrowthResult: skillGrowth,
		}, nil
	}

	return post
}
