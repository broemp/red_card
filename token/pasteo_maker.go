package token

import (
	"errors"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	db "github.com/broemp/red_card/db/sqlc"
	"golang.org/x/crypto/chacha20poly1305"
)

var ErrNotImplementedYet error = errors.New("not implemented yet")

type PasteoMaker struct {
	paseto       paseto.Version
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: mist be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasteoMaker{
		paseto:       "V4",
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (m *PasteoMaker) CreateToken(user db.User, duration time.Duration) (string, *Claims, error) {
	return "", nil, ErrNotImplementedYet
}

func (m *PasteoMaker) VerifyToken(token string) (*Claims, error) {
	return nil, ErrNotImplementedYet
}
