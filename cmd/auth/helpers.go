package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"medods/internal/auth"
	"medods/internal/data"
	"net/http"
	"strings"
	"time"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (tokens *tokens) generate(expires int, secret string) error {
	/*
		Сейчас я использую struct с фиксированными полями для генерации токенов,
		но в реальном приложении можно передавать любой интерфейс
	*/
	accessToken, err := auth.GenerateJWT(struct {
		Name   string
		Access string
	}{
		Name:   "Nikita",
		Access: "admin",
	}, expires, secret)
	if err != nil {
		return err
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return err
	}

	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return nil
}

func (app *application) generateAndInsertToken(guid string) (tokens, error) {
	var tokens tokens

	err := tokens.generate(app.config.token.accessExpires, app.config.token.secret)
	if err != nil {
		return tokens, err
	}

	token := data.Token{
		GUID:          guid,
		RefreshToken:  tokens.RefreshToken,
		RefreshExpiry: time.Now().Add(time.Duration(app.config.token.refreshExpires) * time.Hour),
	}

	err = app.models.Tokens.Insert(&token)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}
