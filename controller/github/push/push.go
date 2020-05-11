package pushcontroller

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

// GithubPush is a struct to which payload is assigned
type GithubPush struct {
	ID             primitive.ObjectID
	ProviderName   string
	Type           string
	ActorName      string
	ActorID        int
	RepositoryName string
	CommitID       int
	CommitMessage  string
	CommitLink     string
	UpdatedTime    time.Time
}

// PayloadBody is struct of the payload
type PayloadBody struct {
	Sender        PayloadSender   `json:"sender"`
	Commits       []PayloadCommit `json:"commits"`
	Repo          PayloadRepo     `json:"repository"`
	DefaultBranch string          `json:"default_branch"`
}

// PayloadSender is struct having sender details fields
type PayloadSender struct {
	UserID   int    `json:"id"`
	UserName string `json:"login"`
}

// PayloadCommit is struct having commit details fields
type PayloadCommit struct {
	CommitID      int    `json:"id"`
	CommitMessage string `json:"message"`
	CommitLink    string `json:"url"`
}

// PayloadRepo is struct having repo details fields
type PayloadRepo struct {
	RepoName string `json:"name"`
}

type pushEventHandler interface {
	push()
}

// PushEvent is used to get payload and assign to json
func PushEvent(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Body)
	var p PayloadBody
	var g GithubPush

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.ProviderName = "github"
	g.Type = "push"
	g.ActorID = p.Sender.UserID
	g.ActorName = p.Sender.UserName
	g.RepositoryName = p.Repo.RepoName
	// g.BranchName = p.Project.Branch
	g.CommitID = p.Commits[0].CommitID
	g.CommitMessage = p.Commits[0].CommitMessage
	g.CommitLink = p.Commits[0].CommitLink
	g.UpdatedTime = time.Now()

	w.Header().Set("Content-Type", "appliction/json")
	json.NewEncoder(w).Encode(g.push())

}

func (g GithubPush) push() *mongo.InsertOneResult {

	db := connection.DBConnection()

	pushEventCollection := db.Database("webhookdb").Collection("githubPush")

	pushResult, err := pushEventCollection.InsertOne(context.TODO(), g)
	if err != nil {
		panic(err)
	}
	return pushResult
}
