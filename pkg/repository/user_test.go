package repository

import (
	"ShowTimes/pkg/utils/models"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCheckUserAvailability(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		stub func(mock sqlmock.Sqlmock)
		want bool
	}{
		{
			name: "successful, user available",
			arg:  "Arun@gmail.com",
			stub: func(mock sqlmock.Sqlmock) {
				querry := "SELECT count(*) FROM USERS WHERE email='Arun@gmail.com'"
				mock.ExpectQuery(regexp.QuoteMeta(querry)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			want: true,
		}, {
			name: "failed, user not avilable",
			arg:  "arun1@gmail.com",
			stub: func(mock sqlmock.Sqlmock) {
				querry := "SELECT count(*) FROM USERS WHERE email='arun1@gmail.com'"
				mock.ExpectQuery(regexp.QuoteMeta(querry)).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want: false,
		},
	}

	for _, tt := range tests {
		mockDb, mockSql, _ := sqlmock.New()
		DB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: mockDb,
		}), &gorm.Config{})
		userRepository := NewUserRepository(DB)
		tt.stub(mockSql)

		result := userRepository.CheckUserAvialiablity(tt.arg)
		assert.Equal(t, tt.want, result)
	}
}

//User Signup Testing

func TestUserSignUp(t *testing.T) {
	type args struct {
		input models.UserDetails
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockSQL sqlmock.Sqlmock)
		want       models.UserDetailsResponse
		wantErr    error
	}{
		{
			name: "Successfully user signed up",
			args: args{
				input: models.UserDetails{
					Name:     "Jojo",
					Email:    "jojo@gmail.com",
					Password: "jojo@123",
					Phone:    "9746359523",
				},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `(?i)insert\s+into\s+users\s*\(\s*name\s*,\s*email\s*,\s*password\s*,\s*phone\s*\)\s*values\s*\(\s*\$1\s*,\s*\$2\s*,\s*\$3\s*,\s*\$4\s*\)\s*returning\s+id\s*,\s*name\s*,\s*email\s*,\s*phone`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("Jojo", "jojo@gmail.com", "jojo@123", "9746359523").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
						AddRow(2, "Jojo", "jojo@gmail.com", "9746359523"))
			},
			want: models.UserDetailsResponse{
				Id:    2,
				Name:  "Jojo",
				Email: "jojo@gmail.com",
				Phone: "9746359523",
			},
			wantErr: nil,
		},
		{
			name: "Error signing up user",
			args: args{
				input: models.UserDetails{
					Name:     "Jojo",
					Email:    "existingemail@gmail.com",
					Password: "jojo@123",
					Phone:    "9746359523",
				},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `(?i)insert\s+into\s+users\s*\(\s*name\s*,\s*email\s*,\s*password\s*,\s*phone\s*\)\s*values\s*\(\s*\$1\s*,\s*\$2\s*,\s*\$3\s*,\s*\$4\s*\)\s*returning\s+id\s*,\s*name\s*,\s*email\s*,\s*phone`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("Jojo", "existingemail@gmail.com", "jojo@123", "9746359523").
					WillReturnError(errors.New("email should be unique"))
			},
			want:    models.UserDetailsResponse{},
			wantErr: errors.New("email should be unique"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSql, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			tt.beforeTest(mockSql)
			u := NewUserRepository(gormDB)
			got, err := u.UserSignup(tt.args.input)
			assert.Equal(t, tt.wantErr, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.UserSignUp() = %v, want %v", got, tt.want)
			}
		})
	}
}
