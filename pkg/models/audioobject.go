package models

import "time"

// AudioObject model
type AudioObject struct {
	ID             IDType    `json:"id"`
	CID            string    `json:"cid"`
	Type           string    `json:"@type"`
	Context        string    `json:"@context"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
	ContentURL     *string   `json:"contentURL" db:"content_url"`
	EncodingFormat *string   `json:"encodingFormat" db:"encoding_format"`
}
