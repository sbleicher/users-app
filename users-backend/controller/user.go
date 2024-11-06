package controller

import (
	"database/sql"
	"errors"
	"strings"
	"users-backend/model"
	"users-backend/repo"
)

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserStatusIncorrect = errors.New("user status is incorrect")
	ErrUsernameCollision   = errors.New("username is already in use")

	_ UserController = new(UserControllerImpl)
)

type UserControllerImpl struct {
	repo repo.UserRepo
}

func NewUserController(repo repo.UserRepo) *UserControllerImpl {
	return &UserControllerImpl{
		repo: repo,
	}
}

func updateUserStatus(userStatus string) (string, error) {
	if userStatus == "A" || userStatus == "I" || userStatus == "T" {
		return userStatus, nil
	} else {
		lowerCaseUserStatus := strings.ToLower(userStatus)
		if lowerCaseUserStatus == "active" {
			return "A", nil
		}
		if lowerCaseUserStatus == "inactive" {
			return "I", nil
		}
		if lowerCaseUserStatus == "terminated" {
			return "T", nil
		}
	}

	return "", ErrUserStatusIncorrect
}

func (c *UserControllerImpl) CreateUser(userName, firstName, lastName, email, userStatus, department string) (int, error) {
	us, err := updateUserStatus(userStatus)
	if err != nil {
		return -1, ErrUserStatusIncorrect
	}

	_, err = c.repo.GetByUsername(userName)
	if err == nil {
		return -1, ErrUserAlreadyExists
	}

	m := &model.User{
		UserName:   userName,
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		UserStatus: us,
		Department: sql.NullString{
			String: department,
			Valid:  true,
		},
	}

	return c.repo.Create(m)
}

func (c *UserControllerImpl) GetAllUsers() (*[]model.User, error) {
	return c.repo.GetAll()
}

func (c *UserControllerImpl) GetUser(user_id int) (*model.User, error) {
	return c.repo.GetById(user_id)
}

func (c *UserControllerImpl) UpdateUser(user_id int, userName, firstName, lastName, email, userStatus, department string) (int, error) {
	us, err := updateUserStatus(userStatus)
	if err != nil {
		return -1, ErrUserStatusIncorrect
	}

	u, err := c.repo.GetByUsername(userName)
	if err == nil && u.UserName == userName && u.UserID != user_id {
		return -1, ErrUsernameCollision
	}

	m := &model.User{
		UserID:     user_id,
		UserName:   userName,
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		UserStatus: us,
		Department: sql.NullString{
			String: department,
			Valid:  true,
		},
	}

	return c.repo.Update(m)
}

func (c *UserControllerImpl) DeleteUser(user_id int) error {
	return c.repo.Delete(user_id)
}
