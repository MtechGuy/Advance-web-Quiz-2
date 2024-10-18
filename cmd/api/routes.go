package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {
	router := httprouter.New()

	// Wrap the handlers with the middleware
	router.Handler(http.MethodGet, "/", a.LoggingMiddleware(http.HandlerFunc(LoggingHandler)))
	router.Handler(http.MethodGet, "/auth", a.LoggingMiddleware(a.AuthMiddleware(http.HandlerFunc(AuthHandler))))
	router.Handler(http.MethodGet, "/handlerrors", a.LoggingMiddleware(a.handleErrorsMiddleware(http.HandlerFunc(helloHandler))))
	router.Handler(http.MethodPost, "/content", a.LoggingMiddleware(a.contentTypeMiddleware(http.HandlerFunc(handleRequest))))

	return router
}
