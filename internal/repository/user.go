package repository

import (
	"github.com/go-pg/pg"
	"github.com/sepuka/campaner/internal/domain"
)

type UserRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Get(userId int) (*domain.User, error) {
	var (
		user = &domain.User{Id: userId}
		err  error
	)

	err = r.
		db.
		Select(user)

	return user, err
}
