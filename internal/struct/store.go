package store

import (
	"github.com/google/uuid"
)

// Represents a recipe
type ShortUrls struct {
	Full string `json:"Full"`
	Id   string `json:"id"`
}

func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}
