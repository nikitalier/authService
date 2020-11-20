package provider

import (
	"context"

	"github.com/nikitalier/authService/config"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Provider interface {
	Open()
	GetCon() *mongo.Client
}

type provider struct {
	DB     *mongo.Client
	Host   string
	logger zerolog.Logger
}

func New(db *config.MongoDB, logger zerolog.Logger) Provider {
	return &provider{
		DB:     nil,
		Host:   db.Host,
		logger: logger,
	}
}

func (p *provider) GetCon() *mongo.Client {
	err := p.DB.Connect(context.TODO())
	if err != nil {
		p.logger.Fatal().Msg(err.Error())
	}
	return p.DB
}

func (p *provider) Open() {
	var err error
	p.DB, err = mongo.NewClient(options.Client().ApplyURI(p.Host))
	if err != nil {
		p.logger.Fatal().Msg(err.Error())
	}
}
