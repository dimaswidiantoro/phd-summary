// models.go
package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chapter struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ChapterTitle string             `bson:"chapterTitle" json:"chapterTitle"`
	Subsections  []Subsection       `bson:"subsections" json:"subsections"`
}

type Subsection struct {
	SubsectionTitle string    `bson:"subsectionTitle" json:"subsectionTitle"`
	Findings        []Finding `bson:"findings" json:"findings"`
}

type Finding struct {
	FindingDescription string   `bson:"findingDescription" json:"findingDescription"`
	SupportingAuthors  []string `bson:"supportingAuthors" json:"supportingAuthors"`
}
