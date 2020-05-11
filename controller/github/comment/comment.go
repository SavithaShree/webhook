package commentcontroller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	connection "webhook/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GithubPRComment is  struct use to assign the payload
type GithubPRComment struct {
	ID             primitive.ObjectID
	ProviderName   string
	Type           string
	ActorName      string
	ActorID        int
	RepositoryName string
	Comment        GithubCommentDetails
	CommitID       int
	PRReviewID     int
	UpdatedTime    time.Time
}
type GithubCommitComment struct {
	ID             primitive.ObjectID
	ProviderName   string
	Type           string
	ActorName      string
	ActorID        int
	RepositoryName string
	Comment        GithubCommentDetails
	CommitID       int
	UpdatedTime    time.Time
}

// GithubCommentDetails is struct for comment details
type GithubCommentDetails struct {
	CommentID  int
	CommentURL string
}

// PayloadBody is struct for payload
type PayloadBody struct {
	Comment PayloadComment `json:"comment"`
	Repo    PayloadRepo    `json:"repo"`
}

// PayloadComment is struct for comment details in payload
type PayloadComment struct {
	CommentURL string              `json:url`
	CommentID  int                 `json:id`
	PRReviewID int                 `json:pull_request_review_id`
	CommitID   int                 `json:"commit_id"`
	Actor      PayloadCommentActor `json:"user"`
}

// PayloadRepo is struct for repo details in payload
type PayloadRepo struct {
	RepoName string `json:"name"`
}

// PayloadCommentActor is struct for details of user who gave comments in payload
type PayloadCommentActor struct {
	ActorID   int    `json:"id"`
	ActorName string `json:"login"`
}

// PullRequestEvent is used to get payload and assign to json
func PullRequestCommentEvent(w http.ResponseWriter, r *http.Request) {

	var p PayloadBody
	var g GithubPRComment

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.ProviderName = "github"
	g.Type = "pull_request_review_comment"
	g.ActorID = p.Comment.Actor.ActorID
	g.ActorName = p.Comment.Actor.ActorName
	g.RepositoryName = p.Repo.RepoName
	g.Comment.CommentID = p.Comment.CommentID
	g.Comment.CommentURL = p.Comment.CommentURL
	g.CommitID = p.Comment.CommitID
	g.PRReviewID = p.Comment.PRReviewID
	g.UpdatedTime = time.Now()

	w.Header().Set("Content-Type", "appliction/json")
	json.NewEncoder(w).Encode(g.pull())

}

// CommitEvent is used to get payload and assign to json
func CommitCommentEvent(w http.ResponseWriter, r *http.Request) {

	var p PayloadBody
	var g GithubCommitComment

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.ProviderName = "github"
	g.Type = "pull_request_review_comment"
	g.ActorID = p.Comment.Actor.ActorID
	g.ActorName = p.Comment.Actor.ActorName
	g.RepositoryName = p.Repo.RepoName
	g.Comment.CommentID = p.Comment.CommentID
	g.Comment.CommentURL = p.Comment.CommentURL
	g.CommitID = p.Comment.CommitID
	g.UpdatedTime = time.Now()

	w.Header().Set("Content-Type", "appliction/json")
	json.NewEncoder(w).Encode(g.pull())

}

func (g GithubPRComment) pull() *mongo.InsertOneResult {

	db := connection.DBConnection()

	prCommentEventCollection := db.Database("webhookdb").Collection("githubPush")

	prComment, err := prCommentEventCollection.InsertOne(context.TODO(), g)
	if err != nil {
		panic(err)
	}
	return prComment
}

func (g GithubCommitComment) pull() *mongo.InsertOneResult {

	db := connection.DBConnection()

	commitCommentEventCollection := db.Database("webhookdb").Collection("githubPush")

	commitComment, err := commitCommentEventCollection.InsertOne(context.TODO(), g)
	if err != nil {
		panic(err)
	}
	return commitComment
}
