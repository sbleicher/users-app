package handler

import (
	"fmt"
	"strconv"
	"users-backend/controller"
	"users-backend/model"

	"github.com/labstack/echo/v4"
)

type (
	HttpUserPost struct {
		UserName   string `json:"user_name"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Email      string `json:"email"`
		UserStatus string `json:"user_status"`
		Department string `json:"department,omitempty"`
	}

	HttpUserPut struct {
		UserID     int    `json:"user_id"`
		UserName   string `json:"user_name"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Email      string `json:"email"`
		UserStatus string `json:"user_status"`
		Department string `json:"department,omitempty"`
	}

	HttpUserIdResponse struct {
		UserID int `json:"user_id"`
	}

	HttpUserResponse struct {
		UserID     int    `json:"user_id"`
		UserName   string `json:"user_name"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Email      string `json:"email"`
		UserStatus string `json:"user_status"`
		Department string `json:"department,omitempty"`
	}

	userHttpHandler struct {
		group      *echo.Group
		controller *controller.UserControllerImpl
	}
)

func NewUserHttpHandler(eg *echo.Group, c *controller.UserControllerImpl) *userHttpHandler {
	return &userHttpHandler{
		group:      eg,
		controller: c,
	}
}

func (h *userHttpHandler) RegisterRoutes() {
	h.group.GET("/:user_id", h.GetUser)
	h.group.GET("", h.GetAllUsers)
	h.group.POST("", h.CreateUser)
	h.group.PUT("", h.UpdateUser)
	h.group.DELETE("/:user_id", h.DeleteUser)
}

func (h *userHttpHandler) CreateUser(c echo.Context) error {
	body := HttpUserPost{}
	if err := c.Bind(&body); err != nil {
		return respError(c, 400, "Invalid body", fmt.Sprintf("Invalid body: %v", err))
	}

	newUserID, err := h.controller.CreateUser(body.UserName, body.FirstName, body.LastName, body.Email, body.UserStatus, body.Department)
	if err != nil {
		if err == controller.ErrUserAlreadyExists {
			return respError(c, 400, "User already exists", fmt.Sprintf("user with username %s already exists", body.UserName))
		} else if err == controller.ErrUserStatusIncorrect {
			return respError(c, 400, "Incorrect Status", fmt.Sprintln("Accepted statuses are: Active, A, Inactive, I, Terminated, T"))
		} else {
			return respError(c, 500, "Internal Server Error", fmt.Sprintf("Unexpected error trying to create user %s", body.UserName))
		}
	}

	type HttpUserPostResponse struct {
		UserID int `json:"user_id"`
	}

	return respSuccess(c, 200, "Success", HttpUserPostResponse{UserID: newUserID})
}

func (h *userHttpHandler) GetAllUsers(c echo.Context) error {
	users, err := h.controller.GetAllUsers()
	if err != nil {
		return respError(c, 500, "Internal Server Error", fmt.Sprintln("Unexpected error trying to get all users"))
	}

	var response []HttpUserResponse
	for _, u := range *users {
		response = append(response, NewHttpUserResponse(u))
	}

	return respSuccess(c, 200, "Success", response)
}

func (h *userHttpHandler) GetUser(c echo.Context) error {
	userIdParam := c.Param("user_id")
	user_id, err := strconv.Atoi(userIdParam)
	if err != nil {
		return respError(c, 400, "Invalid user_id", fmt.Sprintf("user_id %q is not a valid user_id as it is not a number", userIdParam))
	}

	user, err := h.controller.GetUser(user_id)
	if err != nil {
		print(err)
		return respError(c, 500, "Internal Server Error", fmt.Sprintf("Unexpected error trying to get user %q: %s", userIdParam, err))
	}

	return respSuccess(c, 200, "Success", NewHttpUserResponse(*user))
}

func (h *userHttpHandler) UpdateUser(c echo.Context) error {
	c.Echo().Logger.Info("Testing logger within handler1")

	body := HttpUserPut{}
	if err := c.Bind(&body); err != nil {
		return respError(c, 400, "Invalid body", fmt.Sprintf("Invalid body: %v", err))
	}

	updatedUserID, err := h.controller.UpdateUser(body.UserID, body.UserName, body.FirstName, body.LastName, body.Email, body.UserStatus, body.Department)
	if err != nil {
		if err == controller.ErrUsernameCollision {
			return respError(c, 400, "User already exists", fmt.Sprintf("User with username %s already exists", body.UserName))
		} else if err == controller.ErrUserStatusIncorrect {
			return respError(c, 400, "Incorrect Status", fmt.Sprintln("Accepted statuses are: Active, A, Inactive, I, Terminated, T"))
		} else {
			return respError(c, 500, "Internal Server Error", fmt.Sprintf("Unexpected error trying to create user %s", body.UserName))
		}
	}

	type HttpUserPutResponse struct {
		UserID int `json:"user_id"`
	}

	return respSuccess(c, 200, "Success", HttpUserPutResponse{UserID: updatedUserID})

}

func (h *userHttpHandler) DeleteUser(c echo.Context) error {
	userIdParam := c.Param("user_id")
	user_id, err := strconv.Atoi(userIdParam)
	if err != nil {
		return respError(c, 400, "Invalid user_id", fmt.Sprintf("user_id %q is not a valid user_id as it is not a number", userIdParam))
	}

	err = h.controller.DeleteUser(user_id)
	if err != nil {
		return respError(c, 500, "Internal Server Error", fmt.Sprintf("Unexpected error trying to delete user %q", userIdParam))
	}

	return respSuccess(c, 200, "Success")
}

func NewHttpUserResponse(user model.User) HttpUserResponse {
	return HttpUserResponse{
		UserID:     user.UserID,
		UserName:   user.UserName,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		UserStatus: user.UserStatus,
		Department: user.Department.String,
	}
}
