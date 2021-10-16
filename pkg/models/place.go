package models

import "time"

// Place model
// TODO - Delete This and Drop Table
type Place struct {
	ID        IDType    `json:"id"`
	CID       string    `json:"cid"`
	Type      string    `json:"@type"`
	Context   string    `json:"@context"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	Name      string    `json:"name"`
}
