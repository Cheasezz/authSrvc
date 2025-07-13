package core

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}
