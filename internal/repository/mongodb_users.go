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

func (r *mongoDBUsersRepository) Create(ctx context.Context, user entity.UserModel) error {
	_, err := r.coll.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return entity.ErrUserAlreadyExists
	}

	return err
}

func (r *mongoDBUsersRepository) FindByID(ctx context.Context, userID primitive.ObjectID) (entity.UserModel, error) {
	result := r.coll.FindOne(ctx, bson.M{
		"_id": userID,
	})

	var user entity.UserModel
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.UserModel{}, entity.ErrUserNotFound
		}

		return entity.UserModel{}, err
	}

	return user, nil
}

func (r *mongoDBUsersRepository) FindByUsername(ctx context.Context, username string) (entity.UserModel, error) {
	result := r.coll.FindOne(ctx, bson.M{
		"username": username,
	})

	var user entity.UserModel
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.UserModel{}, entity.ErrUserNotFound
		}

		return entity.UserModel{}, err
	}

	return user, nil
}

func (r *mongoDBUsersRepository) FindByEmail(ctx context.Context, email string) (entity.UserModel, error) {
	result := r.coll.FindOne(ctx, bson.M{
		"email": email,
	})

	var user entity.UserModel
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.UserModel{}, entity.ErrUserNotFound
		}

		return entity.UserModel{}, err
	}

	return user, nil
}

func (r *mongoDBUsersRepository) FindByLogin(ctx context.Context, login string) (entity.UserModel, error) {
	result := r.coll.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": login},
			{"email": login},
		},
	})

	var user entity.UserModel
	if err := result.Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.UserModel{}, entity.ErrUserNotFound
		}

		return entity.UserModel{}, err
	}

	return user, nil
}

func (r *mongoDBUsersRepository) ChangeEmail(ctx context.Context, schema ChangeEmailSchema) error {
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": schema.UserID}, bson.M{
		"$set": bson.M{"email": schema.NewEmail},
	})
	if result.MatchedCount == 0 {
		return entity.ErrUserNotFound
	}

	return err
}

func (r *mongoDBUsersRepository) ChangePassword(ctx context.Context, schema ChangePasswordSchema) error {
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": schema.UserID}, bson.M{
		"$set": bson.M{"passwordHash": schema.NewPasswordHash},
	})
	if result.MatchedCount == 0 {
		return entity.ErrUserNotFound
	}

	return err
}
