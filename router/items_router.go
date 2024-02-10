package router

import (
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
	internalServerWrapper := func(w http.ResponseWriter) {
		errorOnInternalError(w, handler.InternalServerError{Message: "given path is invalid"})
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
		if pathSplitLength <= 0 {
			internalServerWrapper(w)
			return
		}
		if pathSplit[0] != routeName {
			internalServerWrapper(w)
			return
		}
		if pathSplitLength == 4 {
			if pathSplit[2] == actionPathName {
				getItemActionDetail(w, r)
				return
			}
			internalServerWrapper(w)
			return
		}
		if pathSplitLength == 2 {
			getItemDetail(w, r)
			return
		}
		if pathSplitLength == 1 {
			getItemList(w, r)
			return
		}
		errorOnNotFound(w, handler.PageNotFoundError{Message: path})
	}

	return h
}
