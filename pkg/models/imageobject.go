package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// ImageObject model
type ImageObject struct {
	ID             IDType    `json:"id"`
	CID            string    `json:"cid"`
	Type           string    `json:"@type"`
	Context        string    `json:"@context"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
	ContentURL     *string   `json:"contentURL" db:"content_url"`
	EncodingFormat *string   `json:"encodingFormat" db:"encoding_format"`
}

// Get fetches an ImageObject by ID
func (obj *ImageObject) Get(db *sqlx.DB, args ...interface{}) (*ImageObject, error) {
	query := `
	SELECT imageobject.id, imageobject.cid, imageobject.type, imageobject.context,
	imageobject.created_at, imageobject.updated_at, imageobject.content_url, imageobject.encoding_format
	FROM imageobject
	WHERE id=$1
	`

	imageobject := &ImageObject{}
	err := db.Get(imageobject, query, args...)

	if err != nil {
		// NOTE: Don't change error message unless you change the string
		// comparison in person.go and anywhere else it's used
		return nil, fmt.Errorf("%v", err)
	}

	return imageobject, err
}
