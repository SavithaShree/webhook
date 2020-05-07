package pushController

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

// GithubPush...
type GithubPush struct {
	ID             primitive.ObjectID
	ProviderName   string
	Type           string
	ActorName      string
	ActorID        int
	RepositoryName string
	BranchName     string
	CommitID       int
	CommitMessage  string
	CommitLink     string
	UpdatedTime    time.Time
}

type PayloadBody struct {
	Sender        PayloadSender   `json:"sender"`
	Commits       []PayloadCommit `json:"commits"`
	Repo          PayloadRepo     `json:"repository"`
	DefaultBranch string          `json:"default_branch"`
}

type PayloadSender struct {
	UserID   int    `json:"id"`
	UserName string `json:"login"`
}
type PayloadCommit struct {
	CommitID      int    `json:"sha"`
	CommitMessage string `json:"message"`
	CommitLink    string `json:"url"`
}

type PayloadRepo struct {
	RepoName string `json:"name"`
}

type pushEventHandler interface {
	push()
}

// PushEvent is used to get payload and assign to json
func PushEvent(w http.ResponseWriter, r *http.Request) {

	log.Println(r)
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

}

func (p GithubPush) push() *mongo.InsertOneResult {

	db := connection.DBConnection()

	pushEventCollection := db.Database("webhookdb").Collection("githubPush")

	pushResult, err := pushEventCollection.InsertOne(context.TODO(), p)
	if err != nil {
		panic(err)
	}
	return pushResult
}
