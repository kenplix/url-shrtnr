package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Kenplix/url-shrtnr/internal/entity"
)

type usersRepository struct {
	coll *mongo.Collection
}

func (m *MongoDB) UsersRepository() *usersRepository {
	return &usersRepository{
		coll: m.db.Collection(usersCollection),
	}
}

func (r *usersRepository) Create(ctx context.Context, user entity.User) error {
	_, err := r.coll.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return entity.ErrUserAlreadyExists
	}

	return err
}

func (r *usersRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
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
