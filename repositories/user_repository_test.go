package repositories

import (
	"context"
	"database/sql/driver"
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
	eighteenYearsAgo := models.Date(
		time.Now().UTC().Add(
			time.Hour * (24 * 365 * 18 * -1),
		),
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
			nil,
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
			nil,
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
				regexp.QuoteMeta("SELECT * FROM `users` WHERE `email` = ? AND deleted_at IS NULL ORDER BY `users`.`id` LIMIT 1"),
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

func (s *userRepositoryTestSuite) TestFindAll() {
	userId := uuid.New()

	tests := []struct {
		description     string
		getUsersError   bool
		countUsersError bool
	}{
		{
			description: "Success",
		},
		{
			description:   "Error getting users",
			getUsersError: true,
		},
		{
			description:     "Error counting users",
			countUsersError: true,
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			expectedGetUsersQuery := s.dbmock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `users` WHERE deleted_at IS NULL LIMIT 10"),
			)
			if test.getUsersError {
				expectedGetUsersQuery.WillReturnError(errors.New("error executing query"))
			} else {
				expectedGetUsersQuery.WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "name", "date_of_birth", "email", "password", "address"},
					).AddRow(
						userId,
						"Stephen Curry",
						time.Date(1988, 3, 14, 0, 0, 0, 0, time.UTC),
						"stephen.curry@nba.com",
						"hashedpass",
						"Av. Paulista, 1000. São Paulo - SP",
					),
				)
			}

			if !test.getUsersError {
				expectedCountUsersQuery := s.dbmock.ExpectQuery(
					regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE deleted_at IS NULL"),
				)
				if test.countUsersError {
					expectedCountUsersQuery.WillReturnError(errors.New("error executing count query"))
				} else {
					expectedCountUsersQuery.WillReturnRows(
						sqlmock.NewRows([]string{"count"}).AddRow(1),
					)
				}
			}

			result, totalResults, err := s.userRepository.FindAll(s.ctx, 10, 0)
			if test.getUsersError || test.countUsersError {
				s.Error(err)
				s.Nil(result)
				s.Equal(int64(0), totalResults)
			} else {
				s.NoError(err)
				s.Equal(int64(1), totalResults)
				s.Equal(result[0].Name, "Stephen Curry")
			}
			s.NoError(s.dbmock.ExpectationsWereMet())
		})
	}
}

func (s *userRepositoryTestSuite) TestFindById() {
	userId := uuid.New()

	tests := []struct {
		description         string
		userId              string
		errorInQuery        error
		noResultsFoundError bool
	}{
		{
			description: "Success",
			userId:      userId.String(),
		},
		{
			description:         "No results found",
			userId:              userId.String(),
			noResultsFoundError: true,
		},
		{
			description:  "Error in query",
			userId:       userId.String(),
			errorInQuery: errors.New("error executing query"),
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			s.SetupTest()

			expectedQuery := s.dbmock.ExpectQuery(
				regexp.QuoteMeta("SELECT * FROM `users` WHERE `id` = ? AND deleted_at IS NULL ORDER BY `users`.`id` LIMIT 1"),
			).WithArgs(
				test.userId,
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
							"email@email.com",
							"hashedpass",
							"Av. Paulista, 1000. São Paulo - SP",
						),
				)
			} else {
				expectedQuery.WillReturnError(test.errorInQuery)
			}

			user, err := s.userRepository.FindById(s.ctx, test.userId)
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

func (s *userRepositoryTestSuite) TestUpdateAttributesByUserId() {
	userId := uuid.New()

	tests := []struct {
		description   string
		attributes    models.User
		expectedQuery string
		expectedArgs  []driver.Value
		errorInQuery  error
	}{
		{
			description:   "Success only name",
			attributes:    models.User{Name: "name"},
			expectedQuery: "UPDATE `users` SET `name`=? WHERE `id` = ?",
			expectedArgs:  []driver.Value{"name", userId.String()},
		},
		{
			description: "Success multiple attributes",
			attributes: models.User{
				Name: "name", Email: "hello@world.com",
			},
			expectedQuery: "UPDATE `users` SET `name`=?,`email`=? WHERE `id` = ?",
			expectedArgs: []driver.Value{
				"name", "hello@world.com", userId.String(),
			},
		},
		{
			description:   "Error in query",
			attributes:    models.User{Name: "name"},
			expectedQuery: "UPDATE `users` SET `name`=? WHERE `id` = ?",
			expectedArgs:  []driver.Value{"name", userId.String()},
			errorInQuery:  errors.New("error executing query"),
		},
	}

	for _, test := range tests {
		s.Run(test.description, func() {
			s.dbmock.ExpectBegin()

			expectedQuery := s.dbmock.ExpectExec(
				regexp.QuoteMeta(test.expectedQuery),
			).WithArgs(
				test.expectedArgs...,
			)
			if test.errorInQuery != nil {
				expectedQuery.WillReturnError(test.errorInQuery)
			} else {
				expectedQuery.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			if test.errorInQuery != nil {
				s.dbmock.ExpectRollback()
			} else {
				s.dbmock.ExpectCommit()
			}

			err := s.userRepository.UpdateAttributesByUserId(
				s.ctx, userId.String(), test.attributes,
			)
			if test.errorInQuery != nil {
				s.ErrorContains(err, test.errorInQuery.Error())
			} else {
				s.NoError(err)
			}
		})
	}
}
