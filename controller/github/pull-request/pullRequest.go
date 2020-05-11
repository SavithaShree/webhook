package pullrequestcontroller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	connection "webhook/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GithubPull is  struct use to assign the payload
type GithubPull struct {
	ID             primitive.ObjectID
	ProviderName   string
	Type           string
	ActorName      string
	ActorID        int
	RepositoryName string
	Assignee       []GithubAssignee
	MergeStatus    bool
	UpdatedTime    time.Time
}

// GithubAssignee is struct for assignee details
type GithubAssignee struct {
	AssigneeID   int
	AssigneeName string
}

// PayloadBody is struct for payload
type PayloadBody struct {
	Sender      PayloadSender      `json:"sender"`
	Repo        PayloadRepo        `json:"repo"`
	PullRequest PayloadPullRequest `json:"pull_request"`
}

// PayloadSender is struct for sender details in payload
type PayloadSender struct {
	UserID   int    `json:"id"`
	UserName string `json:"login"`
}

// PayloadRepo is struct for repo details in payload
type PayloadRepo struct {
	RepoName string `json:"name"`
}

// PayloadPullRequest is struct for pull request details in payload
type PayloadPullRequest struct {
	MergeStatus bool              `json:"merged"`
	Assignee    []PayloadAssignee `json:"assignees"`
}

// PayloadAssignee is struct for assignee details in payload
type PayloadAssignee struct {
	AssigneeID   int    `json:"id"`
	AssigneeName string `json:"login"`
}

// PullRequestEvent is used to get payload and assign to json
func PullRequestEvent(w http.ResponseWriter, r *http.Request) {

	var p PayloadBody
	var g GithubPull

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.ProviderName = "github"
	g.Type = "pull_request"
	g.ActorID = p.Sender.UserID
	g.ActorName = p.Sender.UserName
	g.RepositoryName = p.Repo.RepoName
	g.MergeStatus = p.PullRequest.MergeStatus
	g.UpdatedTime = time.Now()

	// if len(p.PullRequest.Assignee) > 0 {
	// 	for i := 0; i < len(p.PullRequest.Assignee); i++ {
	// 		g.Assignee[i].AssigneeID = p.PullRequest.Assignee[i].AssigneeID
	// 		g.Assignee[i].AssigneeName = p.PullRequest.Assignee[i].AssigneeName
	// 	}
	// } else {
	// 	g.Assignee = []GithubAssignee{}
	// }
	log.Println(len(p.PullRequest.Assignee))

	log.Println(g.ProviderName)
	log.Println(g.Type)
	log.Println(g.ActorID)
	log.Println(g.ActorName)
	log.Println(g.RepositoryName)
	log.Println(g.MergeStatus)
	log.Println(g.UpdatedTime)
	log.Println(g.Assignee)

	w.Header().Set("Content-Type", "appliction/json")
	json.NewEncoder(w).Encode(g.pull())

}

func (g GithubPull) pull() *mongo.InsertOneResult {

	db := connection.DBConnection()

	pullRequestEventCollection := db.Database("webhookdb").Collection("githubPush")

	pullResult, err := pullRequestEventCollection.InsertOne(context.TODO(), g)
	if err != nil {
		panic(err)
	}
	return pullResult
}
