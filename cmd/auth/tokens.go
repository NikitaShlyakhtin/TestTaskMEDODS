package main

import "net/http"

type tokens struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func (app *application) generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Read GUID from the request body

	var tokens tokens

	err := tokens.generate(app.tokenConfig.expires, app.tokenConfig.secret)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: Save hashed Refresh token and its expiration date to the database with GUID

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

	err = tokens.generate(app.tokenConfig.expires, app.tokenConfig.secret)
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
