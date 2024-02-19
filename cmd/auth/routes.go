package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	/*
		Из условия задания не ясно нужно ли ограничивать повторное использование
		маршрута для генерации пары Access и Refresh токенов, поэтому сейчас
		повторный доступ к нему не ограничен.
	*/
	router.HandlerFunc(http.MethodPost, "/auth/token", app.generateTokenHandler)

	router.HandlerFunc(http.MethodPost, "/auth/token/refresh", app.refreshTokenHandler)

	return router
}
