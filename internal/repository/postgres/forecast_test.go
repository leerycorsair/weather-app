package postgres

import (
	"fmt"
	"testing"
	"time"
	"weather-app/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ForecastRepositoryTestSuite struct {
	suite.Suite
	db   *sqlx.DB
	mock sqlmock.Sqlmock
	repo *ForecastRepository
}

func (suite *ForecastRepositoryTestSuite) SetupTest() {
	var err error
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.db = sqlx.NewDb(db, "sqlmock")
	suite.mock = mock
	suite.repo = NewForecastRepository(suite.db)
}

func (suite *ForecastRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *ForecastRepositoryTestSuite) TestCreateForecast() {
	forecast := models.Forecast{
		CityId:       1,
		Temp:         20.5,
		Date:         time.Now(),
		ForecastJson: []byte(`{"weather":"sunny"}`),
	}

	suite.mock.ExpectQuery("insert into forecasts").
		WithArgs(forecast.CityId, forecast.Temp, forecast.Date, forecast.ForecastJson).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := suite.repo.CreateForecast(forecast)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, id)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *ForecastRepositoryTestSuite) TestCreateForecastConflict() {
	forecast := models.Forecast{
		CityId:       1,
		Temp:         20.5,
		Date:         time.Now(),
		ForecastJson: []byte(`{"weather":"sunny"}`),
	}

	suite.mock.ExpectQuery("insert into forecasts").
		WithArgs(forecast.CityId, forecast.Temp, forecast.Date, forecast.ForecastJson).
		WillReturnError(fmt.Errorf("conflict error"))

	id, err := suite.repo.CreateForecast(forecast)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, id)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *ForecastRepositoryTestSuite) TestCreateForecastError() {
	forecast := models.Forecast{
		CityId:       1,
		Temp:         20.5,
		Date:         time.Now(),
		ForecastJson: []byte(`{"weather":"sunny"}`),
	}

	suite.mock.ExpectQuery("insert into forecasts").
		WithArgs(forecast.CityId, forecast.Temp, forecast.Date, forecast.ForecastJson).
		WillReturnError(fmt.Errorf("insertion error"))

	id, err := suite.repo.CreateForecast(forecast)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, id)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *ForecastRepositoryTestSuite) TestGetForecasts() {
	forecasts := []models.Forecast{
		{
			Id:           1,
			CityId:       1,
			Temp:         20.5,
			Date:         time.Now(),
			ForecastJson: []byte(`{"weather":"sunny"}`),
		},
		{
			Id:           2,
			CityId:       1,
			Temp:         22.5,
			Date:         time.Now().Add(24 * time.Hour),
			ForecastJson: []byte(`{"weather":"cloudy"}`),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "city_id", "temp", "date", "forecast_json"}).
		AddRow(forecasts[0].Id, forecasts[0].CityId, forecasts[0].Temp, forecasts[0].Date, forecasts[0].ForecastJson).
		AddRow(forecasts[1].Id, forecasts[1].CityId, forecasts[1].Temp, forecasts[1].Date, forecasts[1].ForecastJson)

	suite.mock.ExpectQuery("select id, city_id, temp, date, forecast_json from forecasts where city_id=\\$1").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := suite.repo.GetForecasts(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), forecasts, result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *ForecastRepositoryTestSuite) TestGetForecastsEmpty() {
	rows := sqlmock.NewRows([]string{"id", "city_id", "temp", "date", "forecast_json"})

	suite.mock.ExpectQuery("select id, city_id, temp, date, forecast_json from forecasts where city_id=\\$1").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := suite.repo.GetForecasts(1)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *ForecastRepositoryTestSuite) TestGetForecastsQueryError() {
	suite.mock.ExpectQuery("select id, city_id, temp, date, forecast_json from forecasts where city_id=\\$1").
		WithArgs(1).
		WillReturnError(fmt.Errorf("query error"))

	result, err := suite.repo.GetForecasts(1)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func TestForecastRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ForecastRepositoryTestSuite))
}
