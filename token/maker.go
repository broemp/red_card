package token

import (
	"time"

	db "github.com/broemp/red_card/db/sqlc"
)

// Maker is an interface for managing tokens
type Maker interface {
	// Creates a new Token for a specific username and duration
	CreateToken(user db.User, duration time.Duration) (string, *Claims, error)

	VerifyToken(token string) (*Claims, error)
}
