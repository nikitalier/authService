package repository

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Repository) AddTokenPairs(tokens map[string]string, guid string, uuid string) error {
	ctx := context.Background()

	hashedRT, err := bcrypt.GenerateFromPassword([]byte(tokens["refresh_token"]), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error().Msg(err.Error())
		return err
	}

	session, err := r.startTransaction(ctx)
	if err != nil {
		return err
	}

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		collT := r.db.Database("authService").Collection("AccessToken")
		_, err := collT.InsertOne(sc, bson.D{{"Token", tokens["access_token"]}, {"GUID", guid}, {"UUID", uuid}})
		if err != nil {
			r.logger.Error().Msg(err.Error())
			return err
		}

		collRT := r.db.Database("authService").Collection("RefreshToken")
		_, err = collRT.InsertOne(sc, bson.D{{"TokenHash", hashedRT}, {"GUID", guid}, {"UUID", uuid}})
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

func (r *Repository) DeleteRefreshToken(uuid string) error {
	var ctx = context.Background()

	session, err := r.startTransaction(ctx)
	if err != nil {
		return err
	}

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		coll := r.db.Database("authService").Collection("RefreshToken")
		_, err := coll.DeleteOne(sc, bson.D{{"UUID", uuid}})
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

func (r *Repository) DeleteAccessToken(uuid string) error {
	var ctx = context.Background()

	session, err := r.startTransaction(ctx)
	if err != nil {
		return err
	}

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		coll := r.db.Database("authService").Collection("AccessToken")
		_, err := coll.DeleteOne(sc, bson.D{{"UUID", uuid}})
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
