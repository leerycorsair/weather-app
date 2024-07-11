package userservice

import (
	"errors"
	"testing"
	"weather-app/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockCityService struct {
	mock.Mock
}

func (m *MockCityService) CreateCity(city models.City) (int, error) {
	args := m.Called(city)
	return args.Int(0), args.Error(1)
}

func (m *MockCityService) GetCities() ([]models.City, error) {
	args := m.Called()
	return args.Get(0).([]models.City), args.Error(1)
}

func (m *MockCityService) GetCity(cityId int) (models.City, error) {
	args := m.Called(cityId)
	return args.Get(0).(models.City), args.Error(1)
}

func (m *MockCityService) FetchCityData(cityName string, openWeatherAPIKey string) (models.City, error) {
	args := m.Called(cityName, openWeatherAPIKey)
	return args.Get(0).(models.City), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user models.User) (int, error) {
	args := m.Called(user)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) GetUser(login, password string) (models.User, error) {
	args := m.Called(login, password)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) GetFavorites(userId int) ([]int, error) {
	args := m.Called(userId)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockUserRepository) AddFavorite(userId int, cityId int) (int, error) {
	args := m.Called(userId, cityId)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) DeleteFavorite(userId int, cityId int) error {
	args := m.Called(userId, cityId)
	return args.Error(0)
}

type UserServiceTestSuite struct {
	suite.Suite
	service     *UserService
	mockCitySvc *MockCityService
	mockUserRep *MockUserRepository
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.mockCitySvc = new(MockCityService)
	suite.mockUserRep = new(MockUserRepository)
	suite.service = NewUserService(suite.mockCitySvc, suite.mockUserRep)
}

func (suite *UserServiceTestSuite) TestCreateUser() {
	user := models.User{
		Login:    "testuser",
		Password: "password",
		Email:    "test@example.com",
	}

	hashedPassword := generatePasswordHash(user.Password)
	suite.mockUserRep.On("CreateUser", mock.MatchedBy(func(u models.User) bool {
		return u.Login == user.Login && u.Email == user.Email && u.Password == hashedPassword
	})).Return(1, nil)

	id, err := suite.service.CreateUser(user)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, id)
	suite.mockUserRep.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestCreateUserError() {
	user := models.User{
		Login:    "testuser",
		Password: "password",
		Email:    "test@example.com",
	}

	suite.mockUserRep.On("CreateUser", mock.Anything).Return(0, errors.New("create user error"))

	id, err := suite.service.CreateUser(user)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, id)
	suite.mockUserRep.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGenerateToken() {
	user := models.User{
		Id:       1,
		Login:    "testuser",
		Password: generatePasswordHash("password"),
	}

	suite.mockUserRep.On("GetUser", user.Login, user.Password).Return(user, nil)

	token, err := suite.service.GenerateToken(user.Login, "password")
	assert.NoError(suite.T(), err)

	claims := &tokenClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Id, claims.UserId)
}

func (suite *UserServiceTestSuite) TestGenerateTokenError() {
	suite.mockUserRep.On("GetUser", "wronguser", generatePasswordHash("password")).Return(models.User{}, assert.AnError)

	token, err := suite.service.GenerateToken("wronguser", "password")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "", token)
}

func (suite *UserServiceTestSuite) TestParseToken() {
	user := models.User{
		Id:       1,
		Login:    "testuser",
		Password: generatePasswordHash("password"),
	}

	suite.mockUserRep.On("GetUser", user.Login, user.Password).Return(user, nil)

	token, err := suite.service.GenerateToken(user.Login, "password")
	assert.NoError(suite.T(), err)

	userId, err := suite.service.ParseToken(token)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Id, userId)
}

func (suite *UserServiceTestSuite) TestGetFavorites() {
	userId := 1
	favorites := []int{1, 2, 3}

	suite.mockUserRep.On("GetFavorites", userId).Return(favorites, nil)

	result, err := suite.service.GetFavorites(userId)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), favorites, result)
	suite.mockUserRep.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestAddFavorite() {
	userId := 1
	cityId := 1

	suite.mockUserRep.On("AddFavorite", userId, cityId).Return(1, nil)

	id, err := suite.service.AddFavorite(userId, cityId)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, id)
	suite.mockUserRep.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestAddFavoriteError() {
	userId := 1
	cityId := 1

	suite.mockUserRep.On("AddFavorite", userId, cityId).Return(0, errors.New("add favorite error"))

	id, err := suite.service.AddFavorite(userId, cityId)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, id)
	suite.mockUserRep.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestDeleteFavorite() {
	userId := 1
	cityId := 1

	suite.mockUserRep.On("DeleteFavorite", userId, cityId).Return(nil)

	err := suite.service.DeleteFavorite(userId, cityId)
	assert.NoError(suite.T(), err)
	suite.mockUserRep.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestDeleteFavoriteError() {
	userId := 1
	cityId := 1

	suite.mockUserRep.On("DeleteFavorite", userId, cityId).Return(errors.New("delete favorite error"))

	err := suite.service.DeleteFavorite(userId, cityId)
	assert.Error(suite.T(), err)
	suite.mockUserRep.AssertExpectations(suite.T())
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
