package test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"users-backend/controller"
	"users-backend/handler"
	"users-backend/model"
	"users-backend/repo/mock"

	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("User Handler", ginkgo.Ordered, func() {
	var (
		mockRepo        *mock.UserRepoMock
		userController  *controller.UserControllerImpl
		e               *echo.Echo
		userHttpHandler *handler.UserHttpHandler
		userJSON        = `{"user_name": "johndoe", "first_name": "John", "last_name": "Doe", "email": "johndoe@email.com", "user_status": "A", "department": "IT"}`
		userJSONUpdate  = `{"user_id": 1, "user_name": "johndoe", "first_name": "John", "last_name": "Doe", "email": "johndoe@email.com", "user_status": "A", "department": "IT"}`
	)

	mockUser := model.User{
		UserName:   "johndoe",
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "johndoe@email.com",
		UserStatus: "A",
		Department: sql.NullString{
			String: "IT",
			Valid:  true,
		}}

	mockUserUpdate := model.User{
		UserID:     1,
		UserName:   "johndoe",
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "johndoe@email.com",
		UserStatus: "A",
		Department: sql.NullString{
			String: "IT",
			Valid:  true,
		}}

	ginkgo.BeforeEach(func() {
		mockRepo = mock.NewUserRepoMock()
		userController = controller.NewUserController(mockRepo)
		e = echo.New()
		group := e.Group("/user")
		userHttpHandler = handler.NewUserHttpHandler(group, userController)
		userHttpHandler.RegisterRoutes()
	})

	ginkgo.Describe("CreateUser", func() {
		ginkgo.It("should return 200 OK w/ user id of created user", func() {
			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(userJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			mockRepo.On("GetByUsername", "johndoe").Return(nil, errors.New("error finding user with username"))
			mockRepo.On("Create", &mockUser).Return(1, nil)

			userHttpHandler.CreateUser(ec)

			var res handler.HttpSuccess
			var resData handler.HttpUserPostResponse
			json.Unmarshal(rec.Body.Bytes(), &res)
			jsonData, _ := json.Marshal(res.Data)
			json.Unmarshal(jsonData, &resData)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusCreated))
			gomega.Expect(resData.UserID).Should(gomega.Equal(1))
		})

		ginkgo.It("should return 200 OK when department is left out", func() {
			_json := `{"user_name": "johndoe", "first_name": "John", "last_name": "Doe", "email": "johndoe@email.com", "user_status": "A"}`
			mockUserD := mockUser
			mockUserD.Department = sql.NullString{
				String: "",
				Valid:  false,
			}

			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(_json))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			mockRepo.On("GetByUsername", "johndoe").Return(nil, errors.New("error finding user with username"))
			mockRepo.On("Create", &mockUserD).Return(1, nil)

			userHttpHandler.CreateUser(ec)

			var res handler.HttpSuccess
			var resData handler.HttpUserPostResponse
			json.Unmarshal(rec.Body.Bytes(), &res)
			jsonData, _ := json.Marshal(res.Data)
			json.Unmarshal(jsonData, &resData)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusCreated))
			gomega.Expect(resData.UserID).Should(gomega.Equal(1))
		})

		ginkgo.It("should return 400 Bad request without a username", func() {
			_json := `{"first_name": "John", "last_name": "Doe", "email": "johndoe@email.com", "user_status": "A"}`

			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(_json))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			userHttpHandler.CreateUser(ec)

			var res handler.HttpError
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("Invalid body"))
		})

		ginkgo.It("should return 400 Bad Request if other data is there", func() {
			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(`"key": "value"`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			userHttpHandler.CreateUser(ec)

			var res handler.HttpSuccess
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("Invalid body"))
		})

		ginkgo.It("should return 400 Bad Request if user already exists", func() {
			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(userJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			mockRepo.On("GetByUsername", "johndoe").Return(&mockUser, nil)

			userHttpHandler.CreateUser(ec)

			var res handler.HttpSuccess
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("User already exists"))
		})

		ginkgo.It("should return 400 Bad Request status is not valid", func() {
			_json := `{"user_name": "johndoe", "first_name": "John", "last_name": "Doe", "email": "johndoe@email.com", "user_status": "ABC", "department": "IT"}`
			req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(_json))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			userHttpHandler.CreateUser(ec)

			var res handler.HttpSuccess
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("Incorrect Status"))
		})
	})

	ginkgo.Describe("UpdateUser", func() {

		ginkgo.It("should return 200 OK w/ user id exists", func() {
			req := httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(userJSONUpdate))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			mockRepo.On("GetByUsername", "johndoe").Return(&mockUserUpdate, nil)
			mockRepo.On("Update", &mockUserUpdate).Return(1, nil)

			userHttpHandler.UpdateUser(ec)

			var res handler.HttpSuccess
			var resData handler.HttpUserPostResponse
			json.Unmarshal(rec.Body.Bytes(), &res)
			jsonData, _ := json.Marshal(res.Data)
			json.Unmarshal(jsonData, &resData)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusOK))
			gomega.Expect(resData.UserID).Should(gomega.Equal(1))
		})

		ginkgo.It("should return 200 OK when department is left out", func() {
			_json := `{"user_id": 1, "user_name": "johndoe", "first_name": "John", "last_name": "Doe", "email": "johndoe@email.com", "user_status": "A"}`
			mockUserD := mockUserUpdate
			mockUserD.Department = sql.NullString{
				String: "",
				Valid:  false,
			}

			req := httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(_json))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			mockRepo.On("GetByUsername", "johndoe").Return(&mockUserD, nil)
			mockRepo.On("Update", &mockUserD).Return(1, nil)

			userHttpHandler.UpdateUser(ec)

			var res handler.HttpSuccess
			var resData handler.HttpUserPostResponse
			json.Unmarshal(rec.Body.Bytes(), &res)
			jsonData, _ := json.Marshal(res.Data)
			json.Unmarshal(jsonData, &resData)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusOK))
			gomega.Expect(resData.UserID).Should(gomega.Equal(1))
		})

		ginkgo.It("should return 400 Bad request without a username", func() {
			_json := `{"first_name": "John", "last_name": "Doe", "email": "johndoe@email.com", "user_status": "A"}`

			req := httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(_json))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			userHttpHandler.UpdateUser(ec)

			var res handler.HttpError
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("Invalid body"))
		})

		ginkgo.It("should return 400 Bad Request if other data is there", func() {
			req := httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(`"key": "value"`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			userHttpHandler.CreateUser(ec)

			var res handler.HttpError
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("Invalid body"))
		})

		ginkgo.It("should return 400 Bad Request if user_name already exists under a different user", func() {
			mockUserD := mockUser
			mockUserD.UserID = 10

			req := httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(userJSONUpdate))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			mockRepo.On("GetByUsername", "johndoe").Return(&mockUserD, nil)

			userHttpHandler.UpdateUser(ec)

			var res handler.HttpSuccess
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("User already exists"))
		})

		ginkgo.It("should return 400 Bad Request status is not valid", func() {
			_json := `{"user_id": 1, "user_name": "johndoe", "first_name": "John", "last_name": "Doe", "email": "johndoe@email.com", "user_status": "ABC", "department": "IT"}`
			req := httptest.NewRequest(http.MethodPut, "/user", strings.NewReader(_json))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			userHttpHandler.UpdateUser(ec)

			var res handler.HttpSuccess
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("Incorrect Status"))
		})
	})

	ginkgo.Describe("GetAllUsers", func() {

		ginkgo.It("should return 200 OK w/ users", func() {
			req := httptest.NewRequest(http.MethodGet, "/user", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			mockRepo.On("GetAll").Return(&[]model.User{mockUser}, nil)

			userHttpHandler.GetAllUsers(ec)

			var res handler.HttpSuccess
			var resData []handler.HttpUserResponse
			json.Unmarshal(rec.Body.Bytes(), &res)
			jsonData, _ := json.Marshal(res.Data)
			json.Unmarshal(jsonData, &resData)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusOK))
			gomega.Expect(len(resData)).Should(gomega.Equal(1))
			gomega.Expect(resData[0].UserName).Should(gomega.Equal("johndoe"))
		})

		ginkgo.It("should return 500 server error when users are not present", func() {
			req := httptest.NewRequest(http.MethodGet, "/user", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			mockRepo.On("GetAll").Return(&[]model.User{}, errors.New("errors"))

			userHttpHandler.GetAllUsers(ec)

			var res handler.HttpError
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusInternalServerError))
			gomega.Expect(res.Message).Should(gomega.Equal("Internal Server Error"))
		})
	})

	ginkgo.Describe("GetUser", func() {

		ginkgo.It("should return 200 OK w/ user", func() {
			req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)
			ec.SetPath("/users/:user_id")
			ec.SetParamNames("user_id")
			ec.SetParamValues("1")

			mockRepo.On("GetById", 1).Return(&mockUser, nil)

			userHttpHandler.GetUser(ec)

			var res handler.HttpSuccess
			var resData handler.HttpUserResponse
			json.Unmarshal(rec.Body.Bytes(), &res)
			jsonData, _ := json.Marshal(res.Data)
			json.Unmarshal(jsonData, &resData)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusOK))
			gomega.Expect(resData.UserName).Should(gomega.Equal("johndoe"))
		})

		ginkgo.It("should return 400 Bad Request when given a bad id", func() {
			req := httptest.NewRequest(http.MethodGet, "/user/abc", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)
			ec.SetPath("/users/:user_id")
			ec.SetParamNames("user_id")
			ec.SetParamValues("abc")

			userHttpHandler.GetUser(ec)

			var res handler.HttpError
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("Invalid user_id"))
		})

		ginkgo.It("should return 400 Bad Request when given a bad id", func() {
			req := httptest.NewRequest(http.MethodGet, "/user/abc", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)
			ec.SetPath("/users/:user_id")
			ec.SetParamNames("user_id")
			ec.SetParamValues("abc")

			userHttpHandler.GetUser(ec)

			var res handler.HttpError
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("Invalid user_id"))
		})

		ginkgo.It("should return 404 when there is no user", func() {
			req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)
			ec.SetPath("/users/:user_id")
			ec.SetParamNames("user_id")
			ec.SetParamValues("1")

			mockRepo.On("GetById", 1).Return(&model.User{}, errors.New("error"))

			userHttpHandler.GetUser(ec)

			var res handler.HttpError
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusNotFound))
			gomega.Expect(res.Message).Should(gomega.Equal("User not found"))
		})
	})

	ginkgo.Describe("DeleteUser", func() {

		ginkgo.It("should return 200 OK", func() {
			req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)
			ec.SetPath("/users/:user_id")
			ec.SetParamNames("user_id")
			ec.SetParamValues("1")

			mockRepo.On("Delete", 1).Return(nil)

			userHttpHandler.DeleteUser(ec)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusOK))
		})

		ginkgo.It("should return 400 Bad Request when given a bad id", func() {
			req := httptest.NewRequest(http.MethodDelete, "/user/abc", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)
			ec.SetPath("/users/:user_id")
			ec.SetParamNames("user_id")
			ec.SetParamValues("abc")

			userHttpHandler.DeleteUser(ec)

			var res handler.HttpError
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("Invalid user_id"))
		})

		ginkgo.It("should return 400 Bad Request when given a bad id", func() {
			req := httptest.NewRequest(http.MethodDelete, "/user/abc", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)
			ec.SetPath("/users/:user_id")
			ec.SetParamNames("user_id")
			ec.SetParamValues("abc")

			userHttpHandler.DeleteUser(ec)

			var res handler.HttpError
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusBadRequest))
			gomega.Expect(res.Message).Should(gomega.Equal("Invalid user_id"))
		})

		ginkgo.It("should return 404 when there is no user", func() {
			req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)
			ec.SetPath("/users/:user_id")
			ec.SetParamNames("user_id")
			ec.SetParamValues("1")

			mockRepo.On("Delete", 1).Return(errors.New("error"))

			userHttpHandler.DeleteUser(ec)

			var res handler.HttpError
			json.Unmarshal(rec.Body.Bytes(), &res)

			gomega.Expect(rec.Code).Should(gomega.Equal(http.StatusNotFound))
			gomega.Expect(res.Message).Should(gomega.Equal("User not found"))
		})
	})
})

func TestUserHandler(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "User Handler Suite")
}
