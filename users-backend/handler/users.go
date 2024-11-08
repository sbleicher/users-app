package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"users-backend/controller"
	"users-backend/model"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	HttpUserPost struct {
		UserName   string  `json:"user_name" validate:"required"`
		FirstName  string  `json:"first_name" validate:"required"`
		LastName   string  `json:"last_name" validate:"required"`
		Email      string  `json:"email" validate:"required,email"`
		UserStatus string  `json:"user_status" validate:"required"`
		Department *string `json:"department,omitempty"`
	}

	HttpUserPut struct {
		UserID     int     `json:"user_id" validate:"required"`
		UserName   string  `json:"user_name" validate:"required"`
		FirstName  string  `json:"first_name" validate:"required"`
		LastName   string  `json:"last_name" validate:"required"`
		Email      string  `json:"email" validate:"required,email"`
		UserStatus string  `json:"user_status" validate:"required"`
		Department *string `json:"department,omitempty"`
	}

	HttpUserIdResponse struct {
		UserID int `json:"user_id"`
	}

	HttpUserResponse struct {
		UserID     int     `json:"user_id"`
		UserName   string  `json:"user_name"`
		FirstName  string  `json:"first_name"`
		LastName   string  `json:"last_name"`
		Email      string  `json:"email"`
		UserStatus string  `json:"user_status"`
		Department *string `json:"department,omitempty"`
	}

	UserHttpHandler struct {
		group      *echo.Group
		controller *controller.UserControllerImpl
	}

	HttpUserPostResponse struct {
		UserID int `json:"user_id"`
	}

	HttpUserPutResponse struct {
		UserID int `json:"user_id"`
	}
)

const success = "Success"

func NewUserHttpHandler(eg *echo.Group, c *controller.UserControllerImpl) *UserHttpHandler {
	return &UserHttpHandler{
		group:      eg,
		controller: c,
	}
}

func (h *UserHttpHandler) RegisterRoutes() {
	h.group.GET("/:user_id", h.GetUser)
	h.group.GET("", h.GetAllUsers)
	h.group.POST("", h.CreateUser)
	h.group.PUT("", h.UpdateUser)
	h.group.DELETE("/:user_id", h.DeleteUser)
}

func (h *UserHttpHandler) CreateUser(c echo.Context) error {
	body := HttpUserPost{}
	validate := validator.New()

	if err := c.Bind(&body); err != nil {
		return respError(c, http.StatusBadRequest, "Invalid body", fmt.Sprintf("Invalid body: %v", err))
	}
	if err := validate.Struct(body); err != nil {
		return respError(c, http.StatusBadRequest, "Invalid body", fmt.Sprintf("Invalid body: %v", err))
	}

	newUserID, err := h.controller.CreateUser(body.UserName, body.FirstName, body.LastName, body.Email, body.UserStatus, pointerToString(body.Department))
	if err != nil {
		if err == controller.ErrUserAlreadyExists {
			return respError(c, http.StatusBadRequest, "User already exists", fmt.Sprintf("user with username %s already exists", body.UserName))
		} else if err == controller.ErrUserStatusIncorrect {
			return respError(c, http.StatusBadRequest, "Incorrect Status", fmt.Sprintln("Accepted statuses are: Active, A, Inactive, I, Terminated, T"))
		} else {
			return respError(c, http.StatusInternalServerError, "Internal Server Error", fmt.Sprintf("Unexpected error trying to create user %s", body.UserName))
		}
	}

	return respSuccess(c, http.StatusCreated, success, HttpUserPostResponse{UserID: newUserID})
}

func (h *UserHttpHandler) GetAllUsers(c echo.Context) error {
	users, err := h.controller.GetAllUsers()
	if err != nil {
		return respError(c, http.StatusInternalServerError, "Internal Server Error", fmt.Sprintln("Unexpected error trying to get all users"))
	}

	var response []HttpUserResponse
	for _, u := range *users {
		response = append(response, NewHttpUserResponse(u))
	}

	return respSuccess(c, http.StatusOK, success, response)
}

func (h *UserHttpHandler) GetUser(c echo.Context) error {
	userIdParam := c.Param("user_id")
	user_id, err := strconv.Atoi(userIdParam)
	if err != nil {
		return respError(c, http.StatusBadRequest, "Invalid user_id", fmt.Sprintf("user_id %q is not a valid user_id as it is not a number", userIdParam))
	}

	user, err := h.controller.GetUser(user_id)
	if err != nil {
		return respError(c, http.StatusNotFound, "User not found", fmt.Sprintf("Unexpected error trying to get user %q: %s", userIdParam, err))
	}

	return respSuccess(c, http.StatusOK, success, NewHttpUserResponse(*user))
}

func (h *UserHttpHandler) UpdateUser(c echo.Context) error {
	body := HttpUserPut{}
	validate := validator.New()

	if err := c.Bind(&body); err != nil {
		return respError(c, http.StatusBadRequest, "Invalid body", fmt.Sprintf("Invalid body: %v", err))
	}
	if err := validate.Struct(body); err != nil {
		return respError(c, http.StatusBadRequest, "Invalid body", fmt.Sprintf("Invalid body: %v", err))
	}

	updatedUserID, err := h.controller.UpdateUser(body.UserID, body.UserName, body.FirstName, body.LastName, body.Email, body.UserStatus, pointerToString(body.Department))
	if err != nil {
		if err == controller.ErrUsernameCollision {
			return respError(c, http.StatusBadRequest, "User already exists", fmt.Sprintf("User with username %s already exists", body.UserName))
		} else if err == controller.ErrUserStatusIncorrect {
			return respError(c, http.StatusBadRequest, "Incorrect Status", fmt.Sprintln("Accepted statuses are: Active, A, Inactive, I, Terminated, T"))
		} else {
			return respError(c, http.StatusInternalServerError, "Internal Server Error", fmt.Sprintf("Unexpected error trying to create user %s", body.UserName))
		}
	}

	return respSuccess(c, http.StatusOK, success, HttpUserPutResponse{UserID: updatedUserID})

}

func (h *UserHttpHandler) DeleteUser(c echo.Context) error {
	userIdParam := c.Param("user_id")
	user_id, err := strconv.Atoi(userIdParam)
	if err != nil {
		return respError(c, http.StatusBadRequest, "Invalid user_id", fmt.Sprintf("user_id %q is not a valid user_id as it is not a number", userIdParam))
	}

	err = h.controller.DeleteUser(user_id)
	if err != nil {
		return respError(c, http.StatusNotFound, "User not found", fmt.Sprintf("Unexpected error trying to delete user %q", userIdParam))
	}

	return respSuccess(c, http.StatusOK, success)
}

func NewHttpUserResponse(user model.User) HttpUserResponse {
	return HttpUserResponse{
		UserID:     user.UserID,
		UserName:   user.UserName,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		UserStatus: user.UserStatus,
		Department: nullStringToPointer(user.Department),
	}
}

func pointerToString(dept *string) string {
	if dept == nil {
		return ""
	}

	return *dept
}

func nullStringToPointer(dept sql.NullString) *string {
	if dept.Valid {
		return &dept.String
	}

	return nil
}
