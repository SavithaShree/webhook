package controller

import (
	"encoding/json"
	"net/http"
)

// Welcome is the base route
func Welcome(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Welcome to webhooks")
}
