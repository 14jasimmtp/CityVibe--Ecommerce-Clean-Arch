package usecase

import (
	"errors"
	"testing"

	mock "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/mock/mockRepo"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepo(ctrl)
	UserUseCase := NewUserUsecase(userRepo)

	tests := map[string]struct {
		input   models.UserSignUpDetails
		stub    func(mock.MockUserRepo, models.UserSignUpDetails)
		WantErr error
	}{
		"otp generation failed": {
			input: models.UserSignUpDetails{
				FirstName:       "jasim",
				LastName:        "mtp",
				Email:           "jasimmtp@gmail.com",
				Password:        "jasi1234",
				ConfirmPassword: "jasi1234",
				Phone:           "9496705233",
			},
			stub: func(mu mock.MockUserRepo, user models.UserSignUpDetails) {
				mu.EXPECT().CheckUserExistsEmail(user.Email)
				mu.EXPECT().CheckUserExistsByPhone(user.Phone)
			},
			WantErr: errors.New("error occured generating otp"),
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			tt.stub(*userRepo, tt.input)
			err := UserUseCase.SignUp(tt.input)

			assert.Equal(t, tt.WantErr, err)
		})
	}
}

// func Test_ResetForgottenPassword(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	userRepo := mock.NewMockUserRepo(ctrl)
// 	UserUseCase := NewUserUsecase(userRepo)

// 	tests := map[string]struct {
// 		input   models.ForgotPassword
// 		stub    func(mock.MockUserRepo, models.ForgotPassword)
// 		WantErr error
// 	}{
// 		"failure": {
// 			input: models.ForgotPassword{
// 				Phone:       "9496705233",
// 				OTP: "1243",
// 				NewPassword: "jasi1234",
// 			},
// 			stub: func(mu mock.MockUserRepo, user models.ForgotPassword) {

// 				mu.EXPECT().ChangePassword(user).Times(1).Return(errors.New("invalid otp entered"))
// 			},
// 			WantErr: errors.New("invalid otp entered"),
// 		},
// 	}

// 	for testname, tt := range tests {
// 		t.Run(testname, func(t *testing.T) {
// 			tt.stub(*userRepo, tt.input)
// 			err := UserUseCase.ResetForgottenPassword(tt.input)
// 			assert.Equal(t, tt.WantErr, err)
// 		})
// 	}
// }

// func Test_AddAddress(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	userRepo := mock.NewMockUserRepo(ctrl)
// 	userUsecase := NewUserUsecase(userRepo)

// 	type args struct {
// 		Address models.Address
// 		token   string
// 	}
// 	tests := map[string]struct {
// 		input   args
// 		stub    func(*mock.MockUserRepo, models.Address, string)
// 		want    models.AddressRes
// 		wantErr error
// 	}{
// 		"failure": {
// 			input: args{
// 				Address: models.Address{
// 					Name:      "jasim",
// 					Housename: "mtphouse",
// 					Phone:     "9496705233",
// 					Street:    "maniyat",
// 					City:      "trikaripur",
// 					State:     "kerala",
// 					Pin:       "671310",
// 				},
// 				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6OSwiZW1haWwiOiJqYXNpbW10cDg0QGdtYWlsLmNvbSIsInJvbGUiOiJ1c2VyIiwiaXNzIjoiY2l0eXZpYmUiLCJleHAiOjE3MDcxOTU3MjQsImlhdCI6MTcwNzE5MjEyNH0.YZJ45UXMZmymOlaFgIkkbmv-SnI29OC1eh03gHWxjP0; Path=/; HttpOnly; Expires=Tue, 06 Feb 2024 05:02:04 GMT;",
// 			},
// 			stub: func(mu *mock.MockUserRepo, address models.Address, token string) {
// 				UserId, err := utils.ExtractUserIdFromToken(token)
// 				if err != nil {
// 					fmt.Println(err)
// 				}
// 				mu.EXPECT().AddAddress(gomock.Eq(address), gomock.Eq(UserId)).Return(models.AddressRes{}, nil)
// 			},
// 			want:    models.AddressRes{},
// 			wantErr: nil,
// 		},
// 	}

// 	for testname, tt := range tests {
// 		t.Run(testname, func(t *testing.T) {
// 			tt.stub(userRepo, tt.input.Address, tt.input.token)
// 			result, err := userUsecase.AddAddress(tt.input.Address, tt.input.token)

// 			assert.Equal(t, tt.want, result)
// 			assert.Equal(t, tt.wantErr, err)
// 		})
// 	}
// }
