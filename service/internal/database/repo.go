package database

import (
	"context"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -destination ./mocks/repo_mock.go -package database . Repository

var _ Repository = new(DataAccessObject)

type (
	Repository interface {
		CardDataAccess
		DeckDataAccess
		GroupDataAccess
		ProviderUsersDataAccess
		UserDataAccess
		SessionDataAccess
		WithTransaction(ctx context.Context, callback func(sessionContext mongo.SessionContext) (interface{}, error), txOptions ...*options.TransactionOptions) error
	}
	DataAccessObject struct {
		db *mongo.Database
		*CardDAO
		*DeckDAO
		*GroupDAO
		*ProviderUsersDAO
		*UserDAO
		*SessionDAO
	}
)

func NewRepository(logger zerolog.Logger, db *mongo.Database) *DataAccessObject {
	l := logger.With().Str("module", "Repository").Logger()
	return &DataAccessObject{
		db,
		NewCardDataAccess(db, l),
		NewDeckDataAccess(db, l),
		NewGroupDataAccess(db, l),
		NewProviderUsersDataAccess(db, l),
		NewUserDataAccess(db, l),
		NewSessionDataAccess(db, l),
	}
}

func (d *DataAccessObject) WithTransaction(ctx context.Context, callback func(sessionContext mongo.SessionContext) (interface{}, error), txOptions ...*options.TransactionOptions) error {
	client := d.db.Client()
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, callback, txOptions...)
	if err != nil {
		return err
	}
	return nil
}
