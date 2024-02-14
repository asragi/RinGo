package router

import (
	"fmt"
	"github.com/asragi/RinGo/core"
	"github.com/asragi/RinGo/handler"
	"net/http"
	"strings"
)

func CreateItemsRouteHandler(
	getItemList handler.Handler,
	getItemDetail handler.Handler,
	getItemActionDetail handler.Handler,
	errorOnMethodNotAllowed handler.ReturnResponseOnErrorFunc,
	errorOnInternalError handler.ReturnResponseOnErrorFunc,
	errorOnNotFound handler.ReturnResponseOnErrorFunc,
) handler.Handler {
	routeName := "items"
	actionPathName := "actions"
	internalServerWrapper := func(w http.ResponseWriter, path string) {
		errorOnInternalError(
			w, core.InternalServerError{
				Message: fmt.Sprintf("given path is invalid: %s", path),
			},
		)
	}
	h := func(w http.ResponseWriter, r *http.Request) {
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
				getItemActionDetail(w, r)
				return
			}
			internalServerWrapper(w, path)
			return
		}
		if pathSplitLength == 3 {
			getItemDetail(w, r)
			return
		}
		if pathSplitLength == 2 {
			getItemList(w, r)
			return
		}
		errorOnNotFound(w, handler.PageNotFoundError{Message: path})
	}

	return h
}
