package repository

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kenplix/url-shrtnr/internal/entity"
)

type mongoDBUsersRepository struct {
	coll *mongo.Collection
}

func (m *mongoDB) createUsersRepository(ctx context.Context) error {
	coll := m.db.Collection(usersCollection)

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return errors.Wrap(err, "failed to crete indices")
	}

	m.users = &mongoDBUsersRepository{
		coll: coll,
	}

	return nil
}

func (m *mongoDB) getUsersRepository() UsersRepository {
	return m.users
}

func (r *mongoDBUsersRepository) Create(ctx context.Context, user entity.User) error {
	_, err := r.coll.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return entity.ErrUserAlreadyExists
	}

	return err
}

func (r *mongoDBUsersRepository) FindByID(ctx context.Context, userID primitive.ObjectID) (entity.User, error) {
	result := r.coll.FindOne(ctx, bson.M{
		"_id": userID,
	})

	var user entity.User
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, err
	}

	return user, nil
}

func (r *mongoDBUsersRepository) FindByUsername(ctx context.Context, username string) (entity.User, error) {
	result := r.coll.FindOne(ctx, bson.M{
		"username": username,
	})

	var user entity.User
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, err
	}

	return user, nil
}

func (r *mongoDBUsersRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	result := r.coll.FindOne(ctx, bson.M{
		"email": email,
	})

	var user entity.User
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, err
	}

	return user, nil
}

func (r *mongoDBUsersRepository) FindByLogin(ctx context.Context, login string) (entity.User, error) {
	result := r.coll.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": login},
			{"email": login},
		},
	})

	var user entity.User
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, err
	}

	return user, nil
}