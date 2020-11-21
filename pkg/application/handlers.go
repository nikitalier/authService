package application

import (
	"encoding/json"
	"net/http"

	"github.com/nikitalier/authService/pkg/models"
)

func (app *Application) getPairTokens(w http.ResponseWriter, r *http.Request) {
	guid, ok := r.URL.Query()["guid"]

	if !ok || len(guid[0]) < 1 {
		app.logger.Error().Msg("Missing url parameter")
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	tokens, err := app.svc.GenerateTokenPair(guid[0])
	if err != nil {
		app.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(tokens)
	if err != nil {
		app.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(j)
}

func (app *Application) refreshToken(w http.ResponseWriter, r *http.Request) {
	var rt models.TokenBody

	err := json.NewDecoder(r.Body).Decode(&rt)
	if err != nil {
		app.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens, err := app.svc.RefreshToken(rt.TokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(tokens)
	if err != nil {
		app.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(j)
}

func (app *Application) deleteRefreshToken(w http.ResponseWriter, r *http.Request) {
	var rt models.TokenBody

	err := json.NewDecoder(r.Body).Decode(&rt)
	if err != nil {
		app.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.svc.DeleteRefreshToken(rt.TokenString)
	if err != nil {
		app.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *Application) deleteAllTokens(w http.ResponseWriter, r *http.Request) {
	guid, ok := r.URL.Query()["guid"]

	if !ok || len(guid[0]) < 1 {
		app.logger.Error().Msg("Missing url parameter")
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	err := app.svc.DeleteAllRefreshTokens(guid[0])
	if err != nil {
		app.logger.Error().Msg(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
