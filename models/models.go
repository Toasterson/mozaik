package models

import "time"

type Tile struct {
	Name string
	Text []byte //TODO Markdown Type
	Courses []Tile
	Pictures []string //TODO Image Processing Library
	Videos []string //TODO Video Processing Library
	Author string //TODO User authentication Library
	Status int
	Asessment Asessment
}

const (
	New = iota
	Review = iota
	Public = iota
	Archived = iota
)

//When a User is Done with a Tile he submits his work (Project) for Asessment
type Asessment struct {
	Asessor string //TODO Link to user that does the asessment
	Feedback []byte //TODO Markdown
	Status int //TODO Enumerator of Asessment Status
	Understanding map[string]int
	OpenQuestions []string
}

const (
	Open = iota
	Progress = iota
	WaitingForFeedback = iota
	Done = iota
)

//When a User work on a Tile a Project is Created to save notes and keep track of the Status
type Project struct {
	Status int
	Colaborators []string //TODO Link to user Entity
	Tile Tile
	Notes []byte //TODO Etherpad integration or other Notetaking Format
	Links []string //TODO Only allow links
}

const (
	Working = iota
	InReview = iota
	Finished = iota
)

type Document struct {
	Tile Tile
	Status int //TODO Enumerator
	Visibility int //TODO Enumerator private/public
	Content []byte
	//TODO Access Restrictions to content
}

type User struct {
	ID   int
	Username string

	// Auth
	Email    string
	Password string

	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Lock
	AttemptNumber int64
	AttemptTime   time.Time
	Locked        time.Time

	// Recover
	RecoverToken       string
	RecoverTokenExpiry time.Time
}

