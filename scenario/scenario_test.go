package scenario

import (
	"context"
	"fmt"
)

type executor interface {
	connectAgent
	useToken
	stageInfoHolder
	stageActionSelector
	itemListHolder
	shelvesHolder
}

func commonScenario(ctx context.Context, c executor) error {
	tryTimes := func(f func() error) error {
		var err error
		for i := 0; i < 3; i++ {
			err = f()
			if err == nil {
				return nil
			}
		}
		return err
	}
	if err := tryTimes(func() error { return getResource(ctx, c) }); err != nil {
		return fmt.Errorf("get resource: %w", err)
	}
	if err := tryTimes(func() error { return getMyShelves(ctx, c) }); err != nil {
		return fmt.Errorf("get my shelves: %w", err)
	}
	if err := tryTimes(func() error { return getItemList(ctx, c) }); err != nil {
		return fmt.Errorf("get item list: %w", err)
	}
	if err := tryTimes(func() error { return getStageList(ctx, c) }); err != nil {
		return fmt.Errorf("get stage list: %w", err)
	}
	if err := tryTimes(func() error { return getStageActionDetail(ctx, c) }); err != nil {
		return fmt.Errorf("get stage action detail: %w", err)
	}
	return nil
}
