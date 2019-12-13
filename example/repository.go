package main

import (
	"github.com/gofrs/uuid"
	"github.com/quickaco/xerosdk/auth"
	"golang.org/x/oauth2"
)

type repository struct {
	sessions map[uuid.UUID]*oauth2.Token
}

// NewRepository will build a new auth.Repository based in a memory persistence
func NewRepository() auth.Repository {
	return &repository{
		sessions: make(map[uuid.UUID]*oauth2.Token),
	}
}

func (r *repository) CreateSession(userID uuid.UUID, t *oauth2.Token) error {
	r.sessions[userID] = t
	return nil
}

func (r *repository) UpdateSession(userID uuid.UUID, t *oauth2.Token) error {
	r.sessions[userID] = t
	return nil
}

func (r *repository) GetSession(userID uuid.UUID) (*oauth2.Token, error) {
	return r.sessions[userID], nil
}
