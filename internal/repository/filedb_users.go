package repository

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/samber/lo"

	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kenplix/url-shrtnr/internal/entity"
)

type fileDBUsersRepository struct {
	Users []entity.UserModel `json:"users"`
	path  string
	mux   sync.RWMutex
	db    *fileDB
}

func (f *fileDB) createUsersRepository() error {
	f.users = &fileDBUsersRepository{
		path: filepath.Join(f.dir, usersCollection+".json"),
		db:   f,
	}

	return f.users.load()
}

func (f *fileDB) getUsersRepository() UsersRepository {
	return f.users
}

func (r *fileDBUsersRepository) Create(_ context.Context, user entity.UserModel) error {
	user.ID = primitive.NewObjectID()

	r.mux.Lock()
	r.Users = append(r.Users, user)
	r.mux.Unlock()

	return r.store()
}

func (r *fileDBUsersRepository) FindByID(_ context.Context, userID primitive.ObjectID) (entity.UserModel, error) {
	r.mux.RLock()
	user, found := lo.Find(r.Users, func(user entity.UserModel) bool {
		return user.ID == userID
	})
	r.mux.RUnlock()

	if !found {
		return entity.UserModel{}, entity.ErrUserNotFound
	}

	return user, nil
}

func (r *fileDBUsersRepository) FindByUsername(_ context.Context, username string) (entity.UserModel, error) {
	r.mux.RLock()
	user, found := lo.Find(r.Users, func(user entity.UserModel) bool {
		return user.Username == username
	})
	r.mux.RUnlock()

	if !found {
		return entity.UserModel{}, entity.ErrUserNotFound
	}

	return user, nil
}

func (r *fileDBUsersRepository) FindByEmail(_ context.Context, email string) (entity.UserModel, error) {
	r.mux.RLock()
	user, found := lo.Find(r.Users, func(user entity.UserModel) bool {
		return user.Email == email
	})
	r.mux.RUnlock()

	if !found {
		return entity.UserModel{}, entity.ErrUserNotFound
	}

	return user, nil
}

func (r *fileDBUsersRepository) FindByLogin(_ context.Context, login string) (entity.UserModel, error) {
	r.mux.RLock()
	user, found := lo.Find(r.Users, func(user entity.UserModel) bool {
		return user.Username == login || user.Email == login
	})
	r.mux.RUnlock()

	if !found {
		return entity.UserModel{}, entity.ErrUserNotFound
	}

	return user, nil
}

func (r *fileDBUsersRepository) load() error {
	r.mux.Lock()
	defer r.mux.Unlock()

	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_RDONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err = dec.Decode(r); err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func (r *fileDBUsersRepository) store() error {
	r.mux.Lock()
	defer r.mux.Unlock()

	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")

	return enc.Encode(r)
}
