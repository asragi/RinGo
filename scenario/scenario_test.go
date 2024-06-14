package scenario

import (
	"context"
	"fmt"
	"time"
)

type errorHolder interface {
	storeError(error)
	hasError() bool
}

type executor interface {
	connectAgent
	loginDataProvider
	tokenHolder
	useToken
	userDataHolder
	stageInfoHolder
	stageActionSelector
	itemListHolder
	itemSelector
	itemDetailHolder
	itemActionSupplier
	shelvesHolder
	shelfSelector
	errorHolder
}

func tryTimes(f func() error) error {
	const retryCount = 3
	var err error
	for i := 0; i < retryCount; i++ {
		err = f()
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed after %d retries: %w", retryCount, err)
}

func firstScenario(ctx context.Context, c executor) error {
	if err := tryTimes(func() error { return signUp(ctx, c) }); err != nil {
		return fmt.Errorf("sign up: %w", err)
	}
	if err := tryTimes(func() error { return login(ctx, c) }); err != nil {
		return fmt.Errorf("login: %w", err)
	}
	if err := tryTimes(func() error { return updateUserName(ctx, c) }); err != nil {
		return fmt.Errorf("update user name: %w", err)
	}
	if err := tryTimes(func() error { return updateShopName(ctx, c) }); err != nil {
		return fmt.Errorf("update shop name: %w", err)
	}
	if err := commonScenario(ctx, c); err != nil {
		return err
	}
	if err := tryTimes(func() error { return updateShelfContent(ctx, c) }); err != nil {
		return fmt.Errorf("update shelf content: %w", err)
	}
	if err := tryTimes(func() error { return updateShelfSize(ctx, c) }); err != nil {
		return fmt.Errorf("update shelf size: %w", err)
	}
	return nil
}

func commonScenario(ctx context.Context, c executor) error {
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
	if err := tryTimes(func() error { return postAction(ctx, c) }); err != nil {
		return fmt.Errorf("post action: %w", err)
	}
	if err := tryTimes(func() error { return getItemList(ctx, c) }); err != nil {
		return fmt.Errorf("get item list: %w", err)
	}
	if err := tryTimes(func() error { return getItemDetail(ctx, c) }); err != nil {
		return fmt.Errorf("get item detail: %w", err)
	}
	if err := tryTimes(func() error { return getItemAction(ctx, c) }); err != nil {
		return fmt.Errorf("get item action: %w", err)
	}
	if err := tryTimes(func() error { return getDailyRanking(ctx, c) }); err != nil {
		return fmt.Errorf("get daily ranking: %w", err)
	}
	return nil
}

type scenarioResult struct {
	err error
}

func execScenario(
	ctx context.Context,
	c executor,
	result chan<- *scenarioResult,
	completed chan<- struct{},
	scenario func(context.Context, executor) error,
) {
	if c.hasError() {
		completed <- struct{}{}
		return
	}
	err := scenario(ctx, c)
	if err != nil {
		c.storeError(err)
		result <- &scenarioResult{err: err}
		completed <- struct{}{}
		return
	}
	result <- &scenarioResult{err: nil}
	completed <- struct{}{}
}

func do(
	ctx context.Context,
	agents []executor,
	scenario func(context.Context, executor) error,
) []*scenarioResult {
	result := make(chan *scenarioResult, len(agents))
	completed := make(chan struct{}, len(agents))
	for _, agent := range agents {
		go execScenario(ctx, agent, result, completed, scenario)
		time.Sleep(time.Millisecond * 50)
	}

	for i := 0; i < len(agents); i++ {
		<-completed
	}
	close(result)
	results := make([]*scenarioResult, 0)
	for r := range result {
		results = append(results, r)
	}
	return results
}

func parallelScenario(ctx context.Context, agentCount int, address string) []*scenarioResult {
	adminCli := newAdmin(address)
	clients := func() []executor {
		result := make([]executor, agentCount)
		for i := 0; i < agentCount; i++ {
			result[i] = newClient(address)
		}
		return result
	}()
	initialTime := time.Now()
	targetTime := initialTime
	results := do(ctx, clients, firstScenario)
	results = append(results, &scenarioResult{err: adminLogin(ctx, adminCli)})
	targetTime = targetTime.Add(40 * time.Minute)
	results = append(results, &scenarioResult{err: changeTime(ctx, adminCli, targetTime)})
	results = append(results, do(ctx, clients, commonScenario)...)
	targetTime = targetTime.Add(10 * time.Minute)
	results = append(results, &scenarioResult{err: changeTime(ctx, adminCli, targetTime)})
	// autoUpdate
	targetTime = targetTime.Add(25 * time.Minute)
	results = append(results, &scenarioResult{err: changeTime(ctx, adminCli, targetTime)})
	results = append(results, do(ctx, clients, commonScenario)...)
	targetTime = targetTime.Add(25 * time.Minute)
	results = append(results, &scenarioResult{err: changeTime(ctx, adminCli, targetTime)})
	results = append(results, &scenarioResult{err: changePeriod(ctx, adminCli)})

	return results
}
