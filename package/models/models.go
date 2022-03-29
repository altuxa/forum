package models

// Users ...
type Users struct {
	Id             int
	Login          string
	Email          string
	Password       string
	AcceptPassword string
}

// Sessions ...
type Sessions struct {
	UserId int
	Value  string
}

// Posts ...
type Posts struct {
	Id         int
	UserId     int
	Title      string
	Text       string
	Author     string
	Tags       []string
	Comment    []Comments
	Tag        []Tags
	IsAuthor   bool
	Likes      int
	Dislikes   int
	IsOnline   bool
	Number     int
	ChangePost bool
	IsRated    bool
}

type Comments struct {
	Id          int
	UserId      int
	PostId      int
	Text        string
	Author      string
	AuthorOfCom bool
	Change      bool
	Likes       int
	Dislikes    int
}

type Tags struct {
	Id     int
	PostId int
	Value  string
}

type RatingPost struct {
	PostId  int
	UserId  int
	Like    int
	Dislike int
}

type RatingComment struct {
	CommentId int
	UserId    int
	Like      int
	Dislike   int
}

type NotificationComments struct {
	Id           int
	UserId       int
	PostId       int
	PostAuthorId int
	CommentId    int
}

type RateNotification struct {
	Id           int
	UserId       int
	PostId       int
	PostAuthorId int
	Rate         string
	RateAuthor   string
}
