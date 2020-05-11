package github

import (
	"log"
	"net/http"
	commentcontroller "webhook/controller/github/comment"
	pullrequestcontroller "webhook/controller/github/pull-request"
	pushcontroller "webhook/controller/github/push"
)

// WebhookEvent will check the event type and triggers the functionality.
func WebhookEvent(w http.ResponseWriter, r *http.Request) {

	header := r.Header.Get("X-GitHub-Event")
	log.Println(header)
	switch header {
	case "push":
		pushcontroller.PushEvent(w, r)
		break
	case "pull_request":
		pullrequestcontroller.PullRequestEvent(w, r)
		break
	case "commit_comment":
		commentcontroller.CommitCommentEvent(w, r)
		break
	case "pull_request_review_comment":
		commentcontroller.PullRequestCommentEvent(w, r)
		break
	}
}
