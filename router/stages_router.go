package router

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/handler"
	"net/http"
	"strings"
)

func CreateStageRouteHandler(
	getStageList handler.Handler,
	getStageActionDetail handler.Handler,
	errorOnMethodNotAllowed handler.ReturnResponseOnErrorFunc,
	errorOnInternalError handler.ReturnResponseOnErrorFunc,
	errorOnNotFound handler.ReturnResponseOnErrorFunc,
) handler.Handler {
	routeName := "stages"
	actionPathName := "actions"
	internalServerWrapper := func(w http.ResponseWriter, path string) {
		errorOnInternalError(
			w, core.InternalServerError{
				Message: fmt.Sprintf("given path is invalid: %s", path),
			},
		)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method != "GET" {
			errorOnMethodNotAllowed(w, handler.MethodNotAllowedError{Message: method})
			return
		}
		path := r.URL.Path
		pathSplit := strings.Split(path, "/")
		pathSplitLength := len(pathSplit)
		if pathSplitLength <= 1 {
			internalServerWrapper(w, path)
			return
		}
		if pathSplit[1] != routeName {
			internalServerWrapper(w, path)
			return
		}
		if pathSplitLength == 5 {
			if pathSplit[3] == actionPathName {
				getStageActionDetail(w, r)
				return
			}
			internalServerWrapper(w, path)
			return
		}
		if pathSplitLength == 2 || (pathSplitLength == 3 && pathSplit[2] == "") {
			getStageList(w, r)
			return
		}
		errorOnNotFound(w, handler.PageNotFoundError{Message: path})
	}
}
