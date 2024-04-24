package router

import (
	"fmt"
	"net/http"
)

type HandleDataRaw struct {
	SamplePathString string
	Method           Method
	Handler          Handler
}

type HandleData struct {
	SamplePath SamplePath
	Method     Method
	Handler    Handler
}

func CreateHandleData(raw []*HandleDataRaw) ([]*HandleData, error) {
	data := make([]*HandleData, len(raw))
	for i, r := range raw {
		samplePath, err := NewSamplePath(r.SamplePathString)
		if err != nil {
			return nil, fmt.Errorf("create handle data: %w", err)
		}
		data[i] = &HandleData{
			SamplePath: samplePath,
			Method:     r.Method,
			Handler:    r.Handler,
		}
	}
	return data, nil
}

func CreateRouter(
	handleDataSet []*HandleData,
) (Handler, error) {
	pathChecks, err := func(data []*HandleData) ([]*pathChecker, error) {
		checks := make([]*pathChecker, len(data))
		for i, d := range data {
			check, err := createCheckRequest(d.Method, d.SamplePath)
			if err != nil {
				return nil, fmt.Errorf("create router: %w", err)
			}
			checks[i] = &pathChecker{
				checkRequestFunc: check,
				handler:          d.Handler,
			}
		}
		return checks, nil
	}(handleDataSet)
	if err != nil {
		return nil, fmt.Errorf("create router: %w", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		path, err := NewPathData(r.URL.Path)
		if err != nil {
			ErrorOnPageNotFound(w, fmt.Errorf("invalid path: %s", r.Method))
			return
		}
		method, err := NewMethod(r.Method)
		if err != nil {
			ErrorOnPageNotFound(w, fmt.Errorf("invalid method: %s", r.Method))
			return
		}
		for _, check := range pathChecks {
			if check.checkRequestFunc(method, path) {
				check.handler(w, r)
				return
			}
		}
		ErrorOnPageNotFound(w, fmt.Errorf("no handler found for path: %s", path.ToString()))
		return
	}, nil
}
