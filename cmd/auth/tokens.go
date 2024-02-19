package main

import (
	"medods/internal/data"
	"net/http"
)

type tokens struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func (app *application) generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		GUID string `json:"guid"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := app.generateAndInsertToken(input.GUID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tokens": tokens}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	guid, err := app.models.Tokens.ValidateToken(input.RefreshToken)
	if err != nil {
		switch err {
		case data.ErrRecordNotFound:
			app.notFoundResponse(w, r)
		case data.ErrRefreshTokenExpired:
			app.tokenExpiredResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	tokens, err := app.generateAndInsertToken(guid)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tokens": tokens}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
