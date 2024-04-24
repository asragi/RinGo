package router

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request)

type pathChecker struct {
	checkRequestFunc checkRequestFunc
	handler          Handler
}

type UseParam[T any] func(*Path) (T, error)

func CreateUseUserIdParam(samplePath SamplePath) UseParam[core.UserId] {
	return func(path *Path) (core.UserId, error) {
		targetParam := paramExpression("{userId}")
		return createUsePathParam[core.UserId](core.CreateUserId, targetParam)(samplePath, path)
	}
}

func createUsePathParam[T any](
	create func(string) (T, error),
	paramExpression paramExpression,
) func(samplePath SamplePath, path *Path) (T, error) {
	return func(samplePath SamplePath, path *Path) (T, error) {
		handleError := func(err error) (T, error) {
			return *new(T), fmt.Errorf("createUsePathParam: %w", err)
		}
		index, err := samplePath.GetParameterIndex(paramExpression)
		if err != nil {
			return handleError(err)
		}
		paramString, err := path.GetStringByIndex(index)
		if err != nil {
			return handleError(err)
		}
		result, err := create(paramString)
		if err != nil {
			return handleError(err)
		}
		return result, nil
	}
}

type checkRequestFunc func(Method, *Path) bool

func createCheckRequest(expectedMethod Method, expectedPath SamplePath) (checkRequestFunc, error) {
	pathMatch, err := NewPathMatchPattern(expectedPath)
	if err != nil {
		return nil, fmt.Errorf("createCheckRequest: %w", err)
	}
	return func(method Method, path *Path) bool {
		if !method.Is(expectedMethod) {
			return false
		}
		if !pathMatch.Match(path) {
			return false
		}
		return true
	}, nil
}
