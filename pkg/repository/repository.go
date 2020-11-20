package repository

import (
	"context"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db     *mongo.Client
	logger zerolog.Logger
}

func New(pr *mongo.Client, logger zerolog.Logger) *Repository {
	return &Repository{
		db:     pr,
		logger: logger,
	}
}

func (r *Repository) PingDB() error {
	err := r.db.Ping(context.TODO(), nil)
	if err != nil {
		r.logger.Error().Msg(err.Error())
	}

	return err
}
