package repository

import (
	"context"

	"github.com/nikitalier/authService/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Repository) AddRefreshToken(hashedRT []byte, uuid, guid string) error {
	ctx := context.Background()

	session, err := r.startTransaction(ctx)
	if err != nil {
		return err
	}

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		coll := r.db.Database("authService").Collection("RefreshToken")
		_, err = coll.InsertOne(sc, bson.D{{"TokenHash", hashedRT}, {"guid", guid}, {"uuid", uuid}})
		if err != nil {
			r.logger.Error().Msg(err.Error())
			return err
		}

		err = session.CommitTransaction(sc)
		if err != nil {
			r.logger.Error().Msg(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	session.EndSession(ctx)
	return nil
}

func (r *Repository) DeleteRefreshTokenByUUID(uuid string) error {
	var ctx = context.Background()

	session, err := r.startTransaction(ctx)
	if err != nil {
		return err
	}

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		coll := r.db.Database("authService").Collection("RefreshToken")
		_, err := coll.DeleteOne(sc, bson.D{{"uuid", uuid}})
		if err != nil {
			r.logger.Error().Msg(err.Error())
			return err
		}

		err = session.CommitTransaction(sc)
		if err != nil {
			r.logger.Error().Msg(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	session.EndSession(ctx)
	return nil
}

func (r *Repository) FindRefreshTokenByUUID(uuid string) (rt models.RefreshToken, err error) {
	var ctx = context.Background()

	session, err := r.startTransaction(ctx)
	if err != nil {
		return rt, err
	}

	filter := bson.D{{"uuid", uuid}}

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		coll := r.db.Database("authService").Collection("RefreshToken")
		err = coll.FindOne(sc, filter).Decode(&rt)
		if err != nil {
			r.logger.Error().Msg(err.Error())
			return err
		}
		err = session.CommitTransaction(sc)
		if err != nil {
			r.logger.Error().Msg(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		return rt, err
	}

	session.EndSession(ctx)
	return rt, err
}

func (r *Repository) DeleteAllRefreshTokensByGUID(guid string) error {
	var ctx = context.Background()

	session, err := r.startTransaction(ctx)
	if err != nil {
		return err
	}

	filter := bson.D{{"guid", guid}}

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		coll := r.db.Database("authService").Collection("RefreshToken")
		_, err := coll.DeleteMany(sc, filter)
		if err != nil {
			r.logger.Error().Msg(err.Error())
			return err
		}

		err = session.CommitTransaction(sc)
		if err != nil {
			r.logger.Error().Msg(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	session.EndSession(ctx)
	return err
}

func (r *Repository) startTransaction(ctx context.Context) (mongo.Session, error) {
	session, err := r.db.StartSession()
	if err != nil {
		r.logger.Error().Msg(err.Error())
		return nil, err
	}

	err = session.StartTransaction()
	if err != nil {
		r.logger.Error().Msg(err.Error())
		return nil, err
	}

	return session, nil
}
