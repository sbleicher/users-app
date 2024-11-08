package mock

import (
	"users-backend/model"
	"users-backend/repo"

	"github.com/stretchr/testify/mock"
)

var (
	_ repo.UserRepo = new(UserRepoMock)
)

type UserRepoMock struct {
	mock.Mock
}

func NewUserRepoMock() *UserRepoMock {
	return &UserRepoMock{}
}

func (r *UserRepoMock) GetById(user_id int) (*model.User, error) {
	args := r.Called(user_id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *UserRepoMock) GetByUsername(userName string) (*model.User, error) {
	args := r.Called(userName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *UserRepoMock) GetAll() (*[]model.User, error) {
	args := r.Called()
	return args.Get(0).(*[]model.User), args.Error(1)
}

func (r *UserRepoMock) Create(user *model.User) (int, error) {
	args := r.Called(user)
	return args.Get(0).(int), args.Error(1)
}

func (r *UserRepoMock) Update(user *model.User) (int, error) {
	args := r.Called(user)
	return args.Get(0).(int), args.Error(1)
}

func (r *UserRepoMock) Delete(user_id int) error {
	args := r.Called(user_id)
	return args.Error(0)
}
