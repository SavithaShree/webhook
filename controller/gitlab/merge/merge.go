package mergeController

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	connection "webhook/db"
)

type mergeEventScheme struct {
	Type              string
	ProviderName      string
	ActorName         string
	ActorID           int
	AssigneeName      string
	AssigneeID        int
	RepositoryName    string
	RespositoryLink   string
	CommitID          string
	CommitMessage     string
	CommitLink        string
	SourceBranch      string
	DestinationBranch string
	MergeStatus       string
	UpdatedTime       time.Time
}

type Test struct {
	Type string `json:"object_kind"`
}

// GitlabPayload is the payload for merge.
type GitlabPayload struct {
	Type              string      `json:"object_kind"`
	ActorName         Actor       `json:"user"`
	Details           MergeDetail `json:"object_attributes"`
	RespositoryDetail Respository `json:"repository"`
}

type Actor struct {
	Name string `json:"username"`
}

type Respository struct {
	Name string `json:"name"`
	Link string `json:"url"`
}

type MergeDetail struct {
	SourceBranch      string   `json:"source_branch"`
	DestinationBranch string   `json:"target_branch"`
	ActorID           int      `json:"author_id"`
	AssgineeID        int      `json:"assignee_id"`
	MergeStatus       string   `json:"merge_status"`
	CommitDetail      Commit   `json:"last_commit"`
	AssigneeDetial    Assignee `json:"assignee"`
}

type Commit struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

type Assignee struct {
	Name string `json:"name"`
}

type mergeEventHandler interface {
	merge()
}

func MergeEvent(w http.ResponseWriter, r *http.Request) {
	var body GitlabPayload
	var p mergeEventScheme

	println(r.Body)
	_ = json.NewDecoder(r.Body).Decode(&body) // Assign the payload to struct.

	p.Type = body.Type
	p.ProviderName = "gitlab"
	p.ActorName = body.ActorName.Name
	p.ActorID = body.Details.ActorID
	p.AssigneeName = body.Details.AssigneeDetial.Name
	p.AssigneeID = body.Details.AssgineeID
	p.RepositoryName = body.RespositoryDetail.Name
	p.RespositoryLink = body.RespositoryDetail.Link
	p.CommitID = body.Details.CommitDetail.ID
	p.CommitMessage = body.Details.CommitDetail.Message
	p.CommitLink = body.Details.CommitDetail.URL
	p.SourceBranch = body.Details.SourceBranch
	p.DestinationBranch = body.Details.DestinationBranch
	p.MergeStatus = body.Details.MergeStatus
	p.UpdatedTime = time.Now()

	w.Header().Set("Content-Type", "appliction/json")
	json.NewEncoder(w).Encode(p.merge())
}

func (m mergeEventScheme) merge() *mongo.InsertOneResult {
	db := connection.DBConnection()

	mergeEventCollection := db.Database("webhookdb").Collection("mergeCollection")

	result, err := mergeEventCollection.InsertOne(context.TODO(), m)
	if err != nil {
		panic(err)
	}
	return result
}
