package store

import (
	"github.com/google/uuid"
)

// Represents a recipe
type ShortUrls struct {
	Full string `form:"Full"`
	Id   string `form:"id"`
}

func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}
