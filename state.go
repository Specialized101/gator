package main

import (
	"github.com/Specialized101/gator/internal/config"
	"github.com/Specialized101/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
