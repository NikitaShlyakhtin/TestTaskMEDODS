package main

import "net/http"

type tokens struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func (app *application) generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var tokens tokens

	err := tokens.generate(app.tokenConfig.expires, app.tokenConfig.secret)
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
	var tokens tokens

	tokens.AccessToken = "access"
	tokens.RefreshToken = "refresh"

	err := app.writeJSON(w, http.StatusOK, envelope{"tokens": tokens}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
