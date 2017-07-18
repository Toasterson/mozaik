package main

import (
	"github.com/toasterson/mozaik/models"
)

type Tile struct {
	Name string
	Text []byte //TODO Markdown Type
	Courses []Tile
	Pictures []string //TODO Image Processing Library
	Videos []string //TODO Video Processing Library
	Author models.User //TODO User authentication Library
	Status int
	Asessment Asessment
}

const (
	TileNew = iota
	TileReview = iota
	TilePublic = iota
	TileArchived = iota
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
	AsessmentOpen = iota
	AsessmentProgress = iota
	AsessmentWaitingForFeedback = iota
	AsessmentDone = iota
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
	ProjectWorking = iota
	ProjectInReview = iota
	ProjectFinished = iota
)

type Document struct {
	Tile Tile
	Status int //TODO Enumerator
	Visibility int //TODO Enumerator private/public
	Content []byte
	//TODO Access Restrictions to content
}

