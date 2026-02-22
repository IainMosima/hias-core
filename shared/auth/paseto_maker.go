package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (TokenMaker, error) {
	if len(symmetricKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       &paseto.V2{},
		symmetricKey: []byte(symmetricKey)[:chacha20poly1305.KeySize],
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(userID string, email string, role string, permissions []string, duration interface{}) (string, *Payload, error) {
	var dur time.Duration
	switch v := duration.(type) {
	case time.Duration:
		dur = v
	default:
		return "", nil, fmt.Errorf("invalid duration type")
	}

	payload, err := NewPayload(userID, email, role, permissions, dur)
	if err != nil {
		return "", payload, err
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, jsonPayload, nil)
	return token, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var jsonPayload []byte
	err := maker.paseto.Decrypt(token, maker.symmetricKey, &jsonPayload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload := &Payload{}
	err = json.Unmarshal(jsonPayload, payload)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
