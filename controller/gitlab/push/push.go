package pushController

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	connection "webhook/db"
)

type pushEventScheme struct {
	Type           string
	ProviderName   string
	ActorName      string
	ActorID        int
	RepositoryName string
	BranchName     string
	CommitID       string
	CommitMessage  string
	CommitLink     string
	UpdatedTime    time.Time
}

type GitlabPayload struct {
	Type      string         `json:"object_kind"`
	UseID     int            `json:"user_id"`
	UserName  string         `json:"user_name"`
	UserEmail string         `json:"user_email"`
	ProjectID int            `json:"project_id"`
	Project   GitlabProject  `json:"project"`
	Commits   []GitlabCommit `json:"commits"`
}

type GitlabProject struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path_with_namespace"`
	Branch string `json:"default_branch"`
}

type GitlabCommit struct {
	ID      string `json:"id"`
	Message string `json:"title"`
	URL     string `json:"url"`
}

type pushEventHandler interface {
	push()
}

func PushEvent(w http.ResponseWriter, r *http.Request) {

	var g GitlabPayload
	var p pushEventScheme

	_ = json.NewDecoder(r.Body).Decode(&g) // Assign the payload to struct.

	p.ProviderName = "gitlab"
	p.Type = g.Type
	p.ActorID = g.UseID
	p.ActorName = g.UserName
	p.RepositoryName = g.Project.Name
	p.BranchName = g.Project.Branch
	p.CommitID = g.Commits[0].ID
	p.CommitMessage = g.Commits[0].Message
	p.CommitLink = g.Commits[0].URL
	p.UpdatedTime = time.Now()

	w.Header().Set("Content-Type", "appliction/json")
	json.NewEncoder(w).Encode(p.push())
}

func (p pushEventScheme) push() *mongo.InsertOneResult {

	db := connection.DBConnection()

	pushEventCollection := db.Database("webhookdb").Collection("pushCollection")

	pushResult, err := pushEventCollection.InsertOne(context.TODO(), p)
	if err != nil {
		panic(err)
	}
	return pushResult
}
