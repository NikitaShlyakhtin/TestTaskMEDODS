package main

import "net/http"

func (app *application) generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var tokens struct {
		AccessToken  string `json:"access"`
		RefreshToken string `json:"refresh"`
	}

	tokens.AccessToken = "access"
	tokens.RefreshToken = "refresh"

	err := app.writeJSON(w, http.StatusOK, envelope{"tokens": tokens}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var tokens struct {
		AccessToken  string `json:"access"`
		RefreshToken string `json:"refresh"`
	}

	tokens.AccessToken = "access"
	tokens.RefreshToken = "refresh"

	err := app.writeJSON(w, http.StatusOK, envelope{"tokens": tokens}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
