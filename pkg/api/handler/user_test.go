package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/usecase/mock"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUsecase := mock.NewMockUserUseCase(ctrl)
	userHandler := NewUserHandler(userUsecase)

	tests := map[string]struct {
		inputData                models.UserLoginDetails
		configureUserUseCaseMock func(mock.MockUserUseCase, models.UserLoginDetails)
		expectedJSON             string
		expectedStatus           int
	}{
		"Valid Test": {
			inputData: models.UserLoginDetails{
				Phone:    "9496705233",
				Password: "jasi432",
			},
			configureUserUseCaseMock: func(mu mock.MockUserUseCase, user models.UserLoginDetails) {
				_, err := utils.Validation(user)
				if err != nil {
					fmt.Println("validation failed")
				}
				mu.EXPECT().UserLogin(user).Times(1).Return(&models.UserLoginResponse{
					ID:        9,
					FirstName: "jasim",
					LastName:  "muhamwed MTP",
					Email:     "jasimmtp84@gmail.com",
					Phone:     "9496705233",
				}, "tokenString", nil)
			},
			expectedJSON:   `{"message":"user successfully logged in","user":{"ID":9,"firstname":"jasim","lastname":"muhamwed MTP","email":"jasimmtp84@gmail.com","phone":"9496705233"}}`,
			expectedStatus: http.StatusOK,
		},
		"Invalid Test": {
			inputData: models.UserLoginDetails{
				Phone:    "9496705233",
				Password: "jasi432",
			},
			configureUserUseCaseMock: func(mu mock.MockUserUseCase, user models.UserLoginDetails) {
				_, err := utils.Validation(user)
				if err != nil {
					fmt.Println("validation failed")
				}
				mu.EXPECT().UserLogin(user).Times(1).Return(nil, "", errors.New("authentication failed"))
			},
			expectedJSON:   `{"error":"authentication failed"}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			tt.configureUserUseCaseMock(*userUsecase, tt.inputData)
			router := gin.Default()

			router.POST("/login", userHandler.UserLogin)
			jsonData, err := json.Marshal(tt.inputData)
			assert.NoError(t, err)
			body := bytes.NewBuffer(jsonData)
			req, err := http.NewRequest("POST", "/login", body)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedJSON, w.Body.String())
		})
	}

}

func Test_UserSignUp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	userUsecase := mock.NewMockUserUseCase(ctrl)
	userHandler := NewUserHandler(userUsecase)

	tests := map[string]struct {
		input                    models.UserSignUpDetails
		configureUserUseCaseMock func(mock.MockUserUseCase, models.UserSignUpDetails)
		ExpectedJSON             string
		ExpectedStatus           int
	}{
		"valid signup": {
			input: models.UserSignUpDetails{
				FirstName:       "jasim",
				LastName:        "mtp",
				Email:           "jasimmtp@gmail.com",
				Password:        "jasi1234",
				ConfirmPassword: "jasi1234",
				Phone:           "9496705233",
			},
			configureUserUseCaseMock: func(mu mock.MockUserUseCase, user models.UserSignUpDetails) {
				_, err := utils.Validation(user)
				if err != nil {
					fmt.Println("validation failed")
				}
				mu.EXPECT().SignUp(user).Times(1).Return(nil)
			},
			ExpectedJSON:   `{"message": "Successfully signed up.Enter otp to login."}`,
			ExpectedStatus: http.StatusOK,
		},
		"Invalid signup": {
			input: models.UserSignUpDetails{
				FirstName:       "jasim",
				LastName:        "mtp",
				Email:           "jasimmtp@gmail.com",
				Password:        "jasi1234",
				ConfirmPassword: "jasi1234",
				Phone:           "9496705233",
			},
			configureUserUseCaseMock: func(mu mock.MockUserUseCase, user models.UserSignUpDetails) {
				_, err := utils.Validation(user)
				if err != nil {
					fmt.Println("validation failed")
				}
				mu.EXPECT().SignUp(user).Times(1).Return(errors.New("authentication error"))
			},
			ExpectedJSON:   `{"error":"authentication error"}`,
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			tt.configureUserUseCaseMock(*userUsecase, tt.input)
			router := gin.Default()
			router.POST("/signup", userHandler.UserSignup)
			jsonData, err := json.Marshal(tt.input)
			assert.NoError(t, err)
			body := bytes.NewBuffer(jsonData)
			req, err := http.NewRequest("POST", "/signup", body)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedStatus, w.Code)
			assert.JSONEq(t, tt.ExpectedJSON, w.Body.String())
		})
	}
}
