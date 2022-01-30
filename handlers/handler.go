package handlers

import (
	"github.com/sajjadanwar0/clubhouse-clone/config"
	"github.com/sajjadanwar0/clubhouse-clone/ent"
)

type Handler struct {
	Client *ent.Client
	Config *config.Config
}

func NewHandlers(client *ent.Client, config *config.Config) *Handler {

	return &Handler{
		Client: client,
		Config: config,
	}
}
