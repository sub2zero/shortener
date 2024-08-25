package store

import (
	"github.com/google/uuid"
)

// Represents a recipe
type ShortUrls struct {
	Full  string `json:"name"`
	Short string `json:"short"`
}

func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}
