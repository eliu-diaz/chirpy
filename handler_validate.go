package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func validateChirp() http.Handler {
	type reqBody struct {
		Body string `json:"body"`
	}

	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	type ValidChirp struct {
		Valid bool `json:"valid"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		reqData := reqBody{}
		err := decoder.Decode(&reqData)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		const maxChirpLength = 140
		if len(reqData.Body) > maxChirpLength {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
			return
		}

		cleanedMessage := censorProfaneMessage(reqData.Body)
		respondWithJSON(w, http.StatusOK, response{CleanedBody: cleanedMessage})
	})
}

func censorProfaneMessage(msg string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	splittedMessage := strings.Split(msg, " ")
	newMessage := []string{}

	for _, word := range splittedMessage {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			newMessage = append(newMessage, "****")
		} else {
			newMessage = append(newMessage, word)
		}
	}

	return strings.Join(newMessage, " ")
}
