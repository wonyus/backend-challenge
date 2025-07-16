package memory

import (
	"context"
	"sync"

	"github.com/wonyus/backend-challenge/internal/domain/entities"
	domainErrors "github.com/wonyus/backend-challenge/internal/domain/errors"
	"github.com/wonyus/backend-challenge/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userRepository struct {
	users  map[primitive.ObjectID]*entities.User
	emails map[string]primitive.ObjectID
	mutex  sync.RWMutex
}

func NewUserRepository() repositories.UserRepository {
	return &userRepository{
		users:  make(map[primitive.ObjectID]*entities.User),
		emails: make(map[string]primitive.ObjectID),
	}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if email already exists
	if _, exists := r.emails[user.Email]; exists {
		return domainErrors.ErrUserAlreadyExists
	}

	r.users[user.ID] = user
	r.emails[user.Email] = user.ID
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domainErrors.ErrUserNotFound
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	userID, exists := r.emails[email]
	if !exists {
		return nil, domainErrors.ErrUserNotFound
	}

	user := r.users[userID]
	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]*entities.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	existingUser, exists := r.users[user.ID]
	if !exists {
		return domainErrors.ErrUserNotFound
	}

	// If email changed, update email mapping
	if existingUser.Email != user.Email {
		delete(r.emails, existingUser.Email)
		r.emails[user.Email] = user.ID
	}

	r.users[user.ID] = user
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user, exists := r.users[id]
	if !exists {
		return domainErrors.ErrUserNotFound
	}

	delete(r.users, id)
	delete(r.emails, user.Email)
	return nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return int64(len(r.users)), nil
}
