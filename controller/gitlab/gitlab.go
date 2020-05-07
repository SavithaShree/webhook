package gitlab

import (
	"net/http"
	mergeController "webhook/controller/gitlab/merge"
	pushController "webhook/controller/gitlab/push"
)

// // EventType is to get the event type of payload.
// type EventType struct {
// 	Type string `json:"object_kind"`
// }

// WebhookEvent will check the event type and triggers the functionality.
func WebhookEvent(w http.ResponseWriter, r *http.Request) {

	// var webhook EventType

	// _ = json.NewDecoder(r.Body).Decode(&webhook)

	header := r.Header.Get("X-Gitlab-Event")
	switch header {
	case "Push Hook":
		pushController.PushEvent(w, r)
		break
	case "Merge Request Hook":
		mergeController.MergeEvent(w, r)
		break
	}
}
