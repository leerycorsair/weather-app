package postgres

import (
	"errors"
	"fmt"
	"testing"
	"weather-app/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *sqlx.DB
	mock sqlmock.Sqlmock
	repo *UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	var err error
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.db = sqlx.NewDb(db, "sqlmock")
	suite.mock = mock
	suite.repo = NewUserRepository(suite.db)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *UserRepositoryTestSuite) TestCreateUser() {
	user := models.User{
		Login:    "testuser",
		Password: "password",
		Email:    "testuser@example.com",
	}

	suite.mock.ExpectQuery("insert into users").
		WithArgs(user.Login, user.Password, user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := suite.repo.CreateUser(user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, id)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestCreateUserError() {
	user := models.User{
		Login:    "testuser",
		Password: "password",
		Email:    "testuser@example.com",
	}

	suite.mock.ExpectQuery("insert into users").
		WithArgs(user.Login, user.Password, user.Email).
		WillReturnError(fmt.Errorf("insertion error"))

	id, err := suite.repo.CreateUser(user)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, id)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetUser() {
	user := models.User{
		Id:       1,
		Login:    "testuser",
		Password: "password",
		Email:    "testuser@example.com",
	}

	suite.mock.ExpectQuery("select id, login, password, email from users where login=\\$1 and password=\\$2").
		WithArgs(user.Login, user.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "email"}).
			AddRow(user.Id, user.Login, user.Password, user.Email))

	result, err := suite.repo.GetUser(user.Login, user.Password)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user, result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetUserNotFound() {
	suite.mock.ExpectQuery("select id, login, password, email from users where login=\\$1 and password=\\$2").
		WithArgs("unknownuser", "wrongpassword").
		WillReturnError(fmt.Errorf("sql: no rows in result set"))

	result, err := suite.repo.GetUser("unknownuser", "wrongpassword")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), models.User{}, result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetUserQueryError() {
	suite.mock.ExpectQuery("select id, login, password, email from users where login=\\$1 and password=\\$2").
		WithArgs("testuser", "password").
		WillReturnError(fmt.Errorf("query error"))

	result, err := suite.repo.GetUser("testuser", "password")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), models.User{}, result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetFavorites() {
	favorites := []int{1, 2, 3}

	suite.mock.ExpectQuery("select city_id from users_cities where user_id=\\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"city_id"}).
			AddRow(favorites[0]).
			AddRow(favorites[1]).
			AddRow(favorites[2]))

	result, err := suite.repo.GetFavorites(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), favorites, result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetFavoritesQueryError() {
	suite.mock.ExpectQuery("select city_id from users_cities where user_id=\\$1").
		WithArgs(1).
		WillReturnError(fmt.Errorf("query error"))

	result, err := suite.repo.GetFavorites(1)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestAddFavorite() {
	suite.mock.ExpectQuery("insert into users_cities").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := suite.repo.AddFavorite(1, 1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, id)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestAddFavoriteError() {
	suite.mock.ExpectQuery("insert into users_cities").
		WithArgs(1, 1).
		WillReturnError(fmt.Errorf("insertion error"))

	id, err := suite.repo.AddFavorite(1, 1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, id)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestDeleteFavorite() {
	suite.mock.ExpectExec("delete from users_cities where user_id=\\$1 and city_id=\\$2").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := suite.repo.DeleteFavorite(1, 1)
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestDeleteFavoriteNoRows() {
	suite.mock.ExpectExec("delete from users_cities where user_id=\\$1 and city_id=\\$2").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := suite.repo.DeleteFavorite(1, 1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), errors.New("no rows deleted"), err)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
