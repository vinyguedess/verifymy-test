package repositories

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"
	"verifymy-golang-test/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type userRepositoryTestSuite struct {
	suite.Suite
	ctx            context.Context
	dbconn         *gorm.DB
	dbmock         sqlmock.Sqlmock
	userRepository UserRepository
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(userRepositoryTestSuite))
}

func (s *userRepositoryTestSuite) SetupTest() {
	s.ctx = context.Background()

	conn, dbmock, _ := sqlmock.New()
	dialector := mysql.Dialector{
		Config: &mysql.Config{
			DSN:                       "sqlmock_db_0",
			Conn:                      conn,
			SkipInitializeWithVersion: true,
		},
	}

	dbconn, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		s.FailNow(err.Error())
	}

	s.dbmock = dbmock
	s.dbconn = dbconn

	s.userRepository = NewUserRepository(s.dbconn)
}

func (s *userRepositoryTestSuite) TestCreate() {
	eighteenYearsAgo := time.Now().UTC().Add(
		time.Hour * (24 * 365 * 18 * -1),
	)

	s.Run("Success", func() {
		s.SetupTest()

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO `users`"),
		).WithArgs(
			sqlmock.AnyArg(),
			"name",
			eighteenYearsAgo,
			"email",
			"hashedpass",
			"Av. Paulista, 1000. São Paulo - SP",
		).WillReturnResult(sqlmock.NewResult(1, 1))
		s.dbmock.ExpectCommit()

		user, err := s.userRepository.Create(
			s.ctx,
			models.User{
				Name:        "name",
				DateOfBirth: eighteenYearsAgo,
				Email:       "email",
				Password:    "hashedpass",
				Address:     "Av. Paulista, 1000. São Paulo - SP",
			},
		)

		s.NoError(err)
		s.Equal(user.Name, "name")
		s.NoError(s.dbmock.ExpectationsWereMet())
	})

	s.Run("Error in query", func() {
		s.SetupTest()

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO `users`"),
		).WithArgs(
			sqlmock.AnyArg(),
			"name",
			eighteenYearsAgo,
			"email",
			"hashedpass",
			"Av. Paulista, 1000. São Paulo - SP",
		).WillReturnError(errors.New("error executing query"))
		s.dbmock.ExpectRollback()

		user, err := s.userRepository.Create(
			s.ctx,
			models.User{
				Name:        "name",
				DateOfBirth: eighteenYearsAgo,
				Email:       "email",
				Password:    "hashedpass",
				Address:     "Av. Paulista, 1000. São Paulo - SP",
			},
		)

		s.ErrorContains(err, "error executing query")
		s.Nil(user)
		s.NoError(s.dbmock.ExpectationsWereMet())
	})
}

func (s *userRepositoryTestSuite) TestFindByEmail() {
	tests := []struct {
		description         string
		email               string
		errorInQuery        error
		noResultsFoundError bool
	}{
		{
			description: "Success",
			email:       "email@email.com",
		},
		{
			description:         "No results found",
			email:               "email@email.com",
			noResultsFoundError: true,
		},
		{
			description:  "Error in query",
			email:        "email@email.com",
			errorInQuery: errors.New("error executing query"),
		},
	}

	userId := uuid.New()

	for _, test := range tests {
		s.Run(test.description, func() {
			s.SetupTest()

			expectedQuery := s.dbmock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `users` WHERE `email` = ? ORDER BY `users`.`id` LIMIT 1"),
			).WithArgs(
				test.email,
			)

			if test.noResultsFoundError {
				expectedQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "date_of_birth", "email", "password", "address"}),
				)
			} else if test.errorInQuery == nil {
				expectedQuery.WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "date_of_birth", "email", "password", "address"}).
						AddRow(
							userId,
							"name",
							time.Now().UTC(),
							test.email,
							"hashedpass",
							"Av. Paulista, 1000. São Paulo - SP",
						),
				)
			} else {
				expectedQuery.WillReturnError(test.errorInQuery)
			}

			user, err := s.userRepository.FindByEmail(s.ctx, test.email)
			if test.errorInQuery != nil {
				s.ErrorContains(err, test.errorInQuery.Error())
				s.Nil(user)
			} else if test.noResultsFoundError {
				s.Nil(user)
				s.Nil(err)
			} else {
				s.NoError(err)
				s.Equal(user.Name, "name")
			}
			s.NoError(s.dbmock.ExpectationsWereMet())
		})
	}
}
