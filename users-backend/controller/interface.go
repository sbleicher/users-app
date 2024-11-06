package controller

import "users-backend/model"

type (
	UserController interface {
		CreateUser(userName, firstName, lastName, email, userStatus, department string) (int, error)
		GetUser(user_id int) (*model.User, error)
		GetAllUsers() (*[]model.User, error)
		UpdateUser(user_id int, userName, firstName, lastName, email, userStatus, department string) (int, error)
		DeleteUser(user_id int) error
	}
)
