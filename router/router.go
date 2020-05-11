package router

import (
	"webhook/controller"
	"webhook/controller/github"
	"webhook/controller/gitlab"
	mergeController "webhook/controller/gitlab/merge"

	"github.com/gorilla/mux"
)

// MakeHTTPHandler will handler all the routes
func MakeHTTPHandler() *mux.Router {
	mux := mux.NewRouter()

	mux.HandleFunc("/", controller.Welcome).Methods("GET")
	mux.HandleFunc("/bitbucket", mergeController.MergeEvent).Methods("POST")
	mux.HandleFunc("/gitlab", gitlab.WebhookEvent).Methods("POST")
	mux.HandleFunc("/github", github.WebhookEvent).Methods("POST")
	return mux
}
