package postgres

import (
	"fmt"
	"testing"
	"weather-app/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CityRepositoryTestSuite struct {
	suite.Suite
	db   *sqlx.DB
	mock sqlmock.Sqlmock
	repo *CityRepository
}

func (suite *CityRepositoryTestSuite) SetupTest() {
	var err error
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.db = sqlx.NewDb(db, "sqlmock")
	suite.mock = mock
	suite.repo = NewCityRepository(suite.db)
}

func (suite *CityRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *CityRepositoryTestSuite) TestCreateCity() {
	city := models.City{
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	suite.mock.ExpectQuery("insert into cities").
		WithArgs(city.Name, city.Country, city.Latitude, city.Longitude).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := suite.repo.CreateCity(city)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, id)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CityRepositoryTestSuite) TestCreateCityConflict() {
	city := models.City{
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	suite.mock.ExpectQuery("insert into cities").
		WithArgs(city.Name, city.Country, city.Latitude, city.Longitude).
		WillReturnError(fmt.Errorf("conflict error"))

	id, err := suite.repo.CreateCity(city)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, id)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CityRepositoryTestSuite) TestGetCities() {
	cities := []models.City{
		{
			Id:        1,
			Name:      "London",
			Country:   "GB",
			Latitude:  51.5074,
			Longitude: -0.1278,
		},
		{
			Id:        2,
			Name:      "New York",
			Country:   "US",
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "country", "latitude", "longitude"}).
		AddRow(cities[0].Id, cities[0].Name, cities[0].Country, cities[0].Latitude, cities[0].Longitude).
		AddRow(cities[1].Id, cities[1].Name, cities[1].Country, cities[1].Latitude, cities[1].Longitude)

	suite.mock.ExpectQuery("select id, name, country, latitude, longitude from cities order by name").WillReturnRows(rows)

	result, err := suite.repo.GetCities()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), cities, result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CityRepositoryTestSuite) TestGetCitiesEmpty() {
	rows := sqlmock.NewRows([]string{"id", "name", "country", "latitude", "longitude"})

	suite.mock.ExpectQuery("select id, name, country, latitude, longitude from cities order by name").WillReturnRows(rows)

	result, err := suite.repo.GetCities()
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CityRepositoryTestSuite) TestGetCity() {
	city := models.City{
		Id:        1,
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "country", "latitude", "longitude"}).
		AddRow(city.Id, city.Name, city.Country, city.Latitude, city.Longitude)

	suite.mock.ExpectQuery("select id, name, country, latitude, longitude from cities where id=\\$1").
		WithArgs(city.Id).
		WillReturnRows(rows)

	result, err := suite.repo.GetCity(city.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), city, result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *CityRepositoryTestSuite) TestGetCityNotFound() {
	suite.mock.ExpectQuery("select id, name, country, latitude, longitude from cities where id=\\$1").
		WithArgs(999).
		WillReturnError(fmt.Errorf("sql: no rows in result set"))

	result, err := suite.repo.GetCity(999)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), models.City{}, result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func TestCityRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CityRepositoryTestSuite))
}
