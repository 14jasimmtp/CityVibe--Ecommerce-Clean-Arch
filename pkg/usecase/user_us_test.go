package usecase

import (
	"errors"
	"testing"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	mock "github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repositorymock"
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