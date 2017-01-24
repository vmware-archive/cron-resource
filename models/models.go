package models

import "time"

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type CheckResponse []Version

type Version struct {
	Time time.Time `json:"time"`
}

type Source struct {
	// min hour dom mon dow
	Expression      string `json:"expression"`
	Location        string `json:"location"`
	FireImmediately bool   `json:"fire_immediately"`
}

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type Metadata []MetadataField

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type InResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}
