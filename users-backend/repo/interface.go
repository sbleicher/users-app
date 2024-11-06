package repo

import (
	"users-backend/model"
)

type (
	UserRepo interface {
		GetById(user_id int) (*model.User, error)
		GetByUsername(userName string) (*model.User, error)
		GetAll() (*[]model.User, error)
		Create(user *model.User) (int, error)
		Update(user *model.User) (int, error)
		Delete(user_id int) error
	}
)
