package repository_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/domain"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/models"
	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetUserById(t *testing.T) {
	tests := map[string]struct {
		ID      int
		stub    func(sqlmock.Sqlmock, int)
		want    *models.UserDetailsResponse
		wantErr error
	}{
		"success": {
			ID: 1,
			stub: func(mockSql sqlmock.Sqlmock, ID int) {
				expectedQuery := `SELECT \* FROM users WHERE id = \$1`
				mockSql.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "phone"}).AddRow(1, "jasim", "mtp", "jasim@gmail.com", "9496705233"))
			},
			want:    &models.UserDetailsResponse{ID: 1, Firstname: "jasim", Lastname: "mtp", Email: "jasim@gmail.com", Phone: "9496705233"},
			wantErr: nil,
		},
		"not found ID": {
			ID: 1,
			stub: func(mockSql sqlmock.Sqlmock, ID int) {
				expectedQuery := `SELECT \* FROM users WHERE id = \$1`
				mockSql.ExpectQuery(expectedQuery).WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: sql.ErrNoRows,
		},
		"error": {
			ID: 1,
			stub: func(mockSql sqlmock.Sqlmock, ID int) {
				expectedQuery := `SELECT \* FROM users WHERE id = \$1`
				mockSql.ExpectQuery(expectedQuery).WillReturnError(errors.New("error fetching user"))
			},
			want:    nil,
			wantErr: errors.New("error fetching user"),
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()
			fmt.Println("before")
			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			if err != nil {
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			fmt.Println("afte")

			if gormDB == nil {
				t.Fatal("gormDB is nil")
			}

			tt.stub(mockSQL, tt.ID)
			fmt.Println("before")

			u := repository.NewUserRepo(gormDB)
			fmt.Println("afte")

			if u == nil {
				t.Fatal("repository is nil")
			}
			fmt.Println("before")

			result, err := u.GetUserById(tt.ID)
			fmt.Println("afte")

			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_CheckUserExistsEmail(t *testing.T) {
	tests := map[string]struct {
		inputEmail string
		stub       func(sqlmock.Sqlmock)
		WantUser   *domain.User
		WantErr    error
	}{

		"user Exist": {
			inputEmail: "jasim@gmail.com",
			stub: func(mockedSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT \* FROM "users" WHERE "users"\."email" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT 1`
				mockedSQL.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "phone"}).AddRow(1, "jasim", "mtp", "jasim@gmail.com", "9496705233"))
			},
			WantUser: &domain.User{
				ID:        1,
				Firstname: "jasim",
				Lastname:  "mtp",
				Email:     "jasim@gmail.com",
				Phone:     "9496705233",
			},
			WantErr: nil,
		},
		"user not Exist": {
			inputEmail: "ja@gmail.com",
			stub: func(mockedSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT \* FROM "users" WHERE "users"\."email" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT 1`
				mockedSQL.ExpectQuery(expectedQuery).WithArgs("ja@gmail.com").
				WillReturnError(sql.ErrNoRows)
			},
			WantUser: nil,
			WantErr:  sql.ErrNoRows,
		},
		"error occured": {
			inputEmail: "jasim@gmail.com",
			stub: func(mockedSQL sqlmock.Sqlmock) {
				expectedQuery := `SELECT \* FROM "users" WHERE "users"\."email" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT 1`
				mockedSQL.ExpectQuery(expectedQuery).WillReturnError(errors.New("error"))
			},
			WantUser: nil,
			WantErr:  errors.New("error"),
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})

			tt.stub(mockSQL)

			u := repository.NewUserRepo(gormDB)

			result, err := u.CheckUserExistsEmail(tt.inputEmail)

			assert.Equal(t, tt.WantUser, result)
			assert.Equal(t, tt.WantErr, err)
		})
	}

}
