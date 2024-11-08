package test

import (
	"database/sql"
	"errors"
	"testing"
	"users-backend/controller"
	"users-backend/model"
	"users-backend/repo/mock"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("User Controller", func() {
	var (
		mockRepo       *mock.UserRepoMock
		userController *controller.UserControllerImpl
	)

	mockUser := model.User{
		UserName:   "username",
		FirstName:  "first",
		LastName:   "last",
		Email:      "username@email.com",
		UserStatus: "A",
		Department: sql.NullString{
			String: "dept.",
			Valid:  true,
		}}

	ginkgo.BeforeEach(func() {
		mockRepo = mock.NewUserRepoMock()
		userController = controller.NewUserController(mockRepo)
	})

	ginkgo.Describe("CreateUser / processCreateUpdateUser", func() {
		ginkgo.It("should return user id of created user", func() {
			mockRepo.On("GetByUsername", "username").Return(nil, errors.New("error finding user with username"))
			mockRepo.On("Create", &mockUser).Return(1, nil)

			val, err := userController.CreateUser("username", "first", "last", "username@email.com", "A", "dept.")

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(val).Should(gomega.Equal(1))
		})

		ginkgo.It("should return error when username already exists", func() {
			mockRepo.On("GetByUsername", "username").Return(&model.User{UserName: "username"}, nil)

			_, err := userController.CreateUser("username", "first", "last", "username@email.com", "A", "dept.")

			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err).Should(gomega.Equal(controller.ErrUserAlreadyExists))
		})

		ginkgo.It("should convert status to single char", func() {
			mockRepo.On("GetByUsername", "username").Return(nil, errors.New("error finding user with username"))
			mockRepo.On("Create", &mockUser).Return(1, nil)

			val, err := userController.CreateUser("username", "first", "last", "username@email.com", "Active", "dept.")

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(val).Should(gomega.Equal(1))
		})

		ginkgo.It("should convert department to null", func() {
			mockUserNull := mockUser
			mockUserNull.Department = sql.NullString{String: "", Valid: false}

			mockRepo.On("GetByUsername", "username").Return(nil, errors.New("error finding user with username"))
			mockRepo.On("Create", &mockUserNull).Return(1, nil)

			val, err := userController.CreateUser("username", "first", "last", "username@email.com", "A", "")

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(val).Should(gomega.Equal(1))
		})

		ginkgo.It("should return error when bad status is given", func() {
			_, err := userController.CreateUser("", "", "", "", "Bad Status", "")

			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(err).Should(gomega.Equal(controller.ErrUserStatusIncorrect))
		})
	})

	ginkgo.Describe("UpdateUser", func() {
		ginkgo.It("should return user id of updated user", func() {
			mockUserUpdate := mockUser
			mockUserUpdate.UserID = 10

			mockRepo.On("GetByUsername", "username").Return(nil, errors.New("error finding user with username"))
			mockRepo.On("Update", &mockUserUpdate).Return(10, nil)

			val, err := userController.UpdateUser(10, "username", "first", "last", "username@email.com", "A", "dept.")

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(val).Should(gomega.Equal(10))
		})
	})

	ginkgo.Describe("GetAllUsers", func() {
		ginkgo.It("should return users", func() {
			mockUsers := []model.User{mockUser}

			mockRepo.On("GetAll").Return(&mockUsers, nil)

			val, err := userController.GetAllUsers()

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(val).Should(gomega.Equal(&mockUsers))
		})
	})

	ginkgo.Describe("GetUser", func() {
		ginkgo.It("should return users", func() {
			mockUserId := mockUser
			mockUserId.UserID = 1

			mockRepo.On("GetById", mockUserId.UserID).Return(&mockUserId, nil)

			val, err := userController.GetUser(mockUserId.UserID)

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(val).Should(gomega.Equal(&mockUserId))
		})
	})

	ginkgo.Describe("GetAllUsers", func() {
		ginkgo.It("should return users", func() {
			mockRepo.On("Delete", mockUser.UserID).Return(nil)

			err := userController.DeleteUser(mockUser.UserID)

			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
	})
})

func TestUserController(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "User Controller Suite")
}
