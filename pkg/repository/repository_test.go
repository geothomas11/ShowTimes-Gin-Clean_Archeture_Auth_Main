package repository

// import (
// 	"regexp"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"gopkg.in/go-playground/assert.v1"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func TestCheckUserAvailability(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		arg  string
// 		stub func(mock sqlmock.Sqlmock)
// 		want bool
// 	}{
// 		{
// 			name: "successful, user available",
// 			arg:  "Arun@gmail.com",
// 			stub: func(mock sqlmock.Sqlmock) {
// 				query := "SELECT count(*) FROM USERS WHERE email='Arun@gmail.com'" // Match actual query
// 				mock.ExpectQuery(regexp.QuoteMeta(query)). // Use correct case & format
// 					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "failed, user not available",
// 			arg:  "Arun1@gmail.com",
// 			stub: func(mock sqlmock.Sqlmock) {
// 				query := "SELECT count(*) FROM USERS WHERE email='Arun1@gmail.com'" // Match actual query
// 				mock.ExpectQuery(regexp.QuoteMeta(query)).
// 					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
// 			},
// 			want: false,
// 		},

// 	}

// 	for _, tt := range tests {
// 		mockDb, mockSql, _ := sqlmock.New()
// 		DB, _ := gorm.Open(postgres.New(postgres.Config{
// 			Conn: mockDb,
// 		}), &gorm.Config{})
// 		userRepository := NewUserRepository(DB)
// 		tt.stub(mockSql)

// 		result := userRepository.CheckUserAvialiablity(tt.arg)
// 		assert.Equal(t, tt.want, result)
// 	}
// }
