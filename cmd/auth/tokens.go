package main

import (
	"medods/internal/data"
	"net/http"
	"time"
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

	var tokens tokens

	err = tokens.generate(app.config.token.accessExpires, app.config.token.secret)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token := data.Token{
		GUID:          input.GUID,
		RefreshToken:  tokens.RefreshToken,
		RefreshExpiry: time.Now().Add(time.Duration(app.config.token.refreshExpires) * time.Hour),
	}

	err = app.models.Tokens.Insert(&token)
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

	// TODO: Verify that the Refresh token is valid and not expired

	var tokens tokens

	err = tokens.generate(app.config.token.accessExpires, app.config.token.secret)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: If the Refresh token was valid then update hashed Refresh token in the database with GUID

	err = app.writeJSON(w, http.StatusOK, envelope{"tokens": tokens}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
