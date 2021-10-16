package models

import "time"

// TODO - Deprecate this model
// Copyright model
type Copyright struct {
	ID           IDType      `json:"id"`
	CID          string      `json:"cid"`
	Type         string      `json:"@type"`
	Context      string      `json:"@context"`
	CreatedAt    time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time   `json:"updatedAt" db:"updated_at"`
	RightsOf     interface{} `json:"rightsOf"`
	ValidFrom    string      `json:"validFrom"`
	ValidThrough string      `json:"validThrough"`
}
