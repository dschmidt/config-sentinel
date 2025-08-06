package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type rule struct {
	Path  string "yaml:\"path\""
	Rule  string "yaml:\"rule\""
	Value string "yaml:\"value\""
}

type config struct {
	Rules []rule "yaml:\"rules\""
}

func respondWithError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func inputTransformHandler(s *ServeCmd) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		claims, err := extractClaimsWithoutValidation(authHeader)
		if err != nil {
			respondWithError(w, fmt.Sprintf("Error extracting claims: %v\n", err), http.StatusUnauthorized)
			return
		}

		rulesContent, err := os.ReadFile(s.Rules)
		if err != nil {
			respondWithError(w, fmt.Sprintf("Error reading rules: %v\n", err), http.StatusInternalServerError)
			return
		}

		config := config{}
		if err := yaml.Unmarshal([]byte(rulesContent), &config); err != nil {
			respondWithError(w, fmt.Sprintf("Error parsing rules: %v\n", err), http.StatusInternalServerError)
			return
		}

		content, err := os.ReadFile(s.Input)
		if err != nil {
			respondWithError(w, fmt.Sprintf("Error reading input: %v\n", err), http.StatusInternalServerError)
			return
		}

		data := map[string]any{}
		if err := yaml.Unmarshal(content, &data); err != nil {
			respondWithError(w, fmt.Sprintf("Error parsing input: %v\n", err), http.StatusInternalServerError)
			return
		}

		err = applyRules(data, config.Rules, claims)
		if err != nil {
			respondWithError(w, fmt.Sprintf("Error applying rules on input: %v\n", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			panic(err)
		}
	}
}
