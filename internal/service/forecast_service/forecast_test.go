package forecastservice

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"weather-app/internal/models"

	"github.com/jarcoal/httpmock"
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

type MockForecastRepository struct {
	mock.Mock
}

func (m *MockForecastRepository) CreateForecast(forecast models.Forecast) (int, error) {
	args := m.Called(forecast)
	return args.Int(0), args.Error(1)
}

func (m *MockForecastRepository) GetForecasts(cityId int) ([]models.Forecast, error) {
	args := m.Called(cityId)
	return args.Get(0).([]models.Forecast), args.Error(1)
}

type ForecastServiceTestSuite struct {
	suite.Suite
	service         *ForecastService
	mockCitySvc     *MockCityService
	mockForecastRep *MockForecastRepository
	apiKey          string
}

func (suite *ForecastServiceTestSuite) SetupTest() {
	suite.mockCitySvc = new(MockCityService)
	suite.mockForecastRep = new(MockForecastRepository)
	suite.service = NewForecastService(suite.mockCitySvc, suite.mockForecastRep)
	suite.apiKey = "test-api-key"
	httpmock.Activate()
}

func (suite *ForecastServiceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *ForecastServiceTestSuite) TestCreateForecast() {
	forecast := models.Forecast{
		CityId:       1,
		Temp:         20.5,
		Date:         time.Now(),
		ForecastJson: []byte(`{"weather":"sunny"}`),
	}

	suite.mockForecastRep.On("CreateForecast", forecast).Return(1, nil)

	id, err := suite.service.CreateForecast(forecast)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, id)
	suite.mockForecastRep.AssertExpectations(suite.T())
}

func (suite *ForecastServiceTestSuite) TestCreateForecastError() {
	forecast := models.Forecast{
		CityId:       1,
		Temp:         20.5,
		Date:         time.Now(),
		ForecastJson: []byte(`{"weather":"sunny"}`),
	}

	suite.mockForecastRep.On("CreateForecast", forecast).Return(0, fmt.Errorf("insert error"))

	id, err := suite.service.CreateForecast(forecast)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, id)
	suite.mockForecastRep.AssertExpectations(suite.T())
}

func (suite *ForecastServiceTestSuite) TestGetShortForecast() {
	city := models.City{
		Id:        1,
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	forecasts := []models.Forecast{
		{
			CityId:       1,
			Temp:         20.5,
			Date:         time.Now().Add(24 * time.Hour),
			ForecastJson: []byte(`{"weather":"sunny"}`),
		},
		{
			CityId:       1,
			Temp:         22.5,
			Date:         time.Now().Add(48 * time.Hour),
			ForecastJson: []byte(`{"weather":"cloudy"}`),
		},
	}

	suite.mockCitySvc.On("GetCity", 1).Return(city, nil)
	suite.mockForecastRep.On("GetForecasts", 1).Return(forecasts, nil)

	result, err := suite.service.GetShortForecast(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), city.Name, result.City)
	assert.Equal(suite.T(), city.Country, result.Country)
	assert.Equal(suite.T(), []string{time.Now().Add(24 * time.Hour).Format("2006-01-02"), time.Now().Add(48 * time.Hour).Format("2006-01-02")}, result.AvailableDates)
	assert.Equal(suite.T(), float32(21.5), result.AvgTemp)
	suite.mockCitySvc.AssertExpectations(suite.T())
	suite.mockForecastRep.AssertExpectations(suite.T())
}

func (suite *ForecastServiceTestSuite) TestGetShortForecastError() {
	suite.mockCitySvc.On("GetCity", 1).Return(models.City{}, fmt.Errorf("city not found"))

	result, err := suite.service.GetShortForecast(1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), models.ForecastSummary{}, result)
	suite.mockCitySvc.AssertExpectations(suite.T())
}

func (suite *ForecastServiceTestSuite) TestGetDetailedForecast() {
	cityId := 1
	date := time.Now().Add(24 * time.Hour)

	forecasts := []models.Forecast{
		{
			CityId:       cityId,
			Temp:         20.5,
			Date:         date,
			ForecastJson: []byte(`{"weather":"sunny"}`),
		},
	}

	suite.mockForecastRep.On("GetForecasts", cityId).Return(forecasts, nil)

	result, err := suite.service.GetDetailedForecast(cityId, date)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), forecasts, result)
	suite.mockForecastRep.AssertExpectations(suite.T())
}

func (suite *ForecastServiceTestSuite) TestGetDetailedForecastNoResults() {
	cityId := 1
	date := time.Now().Add(24 * time.Hour)

	suite.mockForecastRep.On("GetForecasts", cityId).Return([]models.Forecast{}, nil)

	result, err := suite.service.GetDetailedForecast(cityId, date)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "no forecasts were found", err.Error())
	suite.mockForecastRep.AssertExpectations(suite.T())
}

func (suite *ForecastServiceTestSuite) TestFetchForecastData() {
	city := models.City{
		Id:        1,
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	forecastResponse := forecastResponse{
		List: []forecastItem{
			{
				Dt:   time.Now().Unix(),
				Main: mainData{Temp: 20.5},
			},
		},
	}
	responseBody, _ := json.Marshal(forecastResponse)
	httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s", city.Latitude, city.Longitude, suite.apiKey),
		httpmock.NewStringResponder(200, string(responseBody)))

	expectedForecast := models.Forecast{
		CityId: city.Id,
		Temp:   20.5,
		Date:   time.Now(),
	}

	result, err := suite.service.FetchForecastData(city, suite.apiKey)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), expectedForecast.CityId, result[0].CityId)
	assert.Equal(suite.T(), expectedForecast.Temp, result[0].Temp)
	assert.WithinDuration(suite.T(), expectedForecast.Date, result[0].Date, time.Second)
	suite.mockCitySvc.AssertExpectations(suite.T())
	suite.mockForecastRep.AssertExpectations(suite.T())
}

func (suite *ForecastServiceTestSuite) TestFetchForecastDataError() {
	city := models.City{
		Id:        1,
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s", city.Latitude, city.Longitude, suite.apiKey),
		httpmock.NewErrorResponder(fmt.Errorf("network error")))

	result, err := suite.service.FetchForecastData(city, suite.apiKey)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *ForecastServiceTestSuite) TestFetchForecastDataInvalidJSON() {
	city := models.City{
		Id:        1,
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s", city.Latitude, city.Longitude, suite.apiKey),
		httpmock.NewStringResponder(200, "{invalid json}"))

	result, err := suite.service.FetchForecastData(city, suite.apiKey)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *ForecastServiceTestSuite) TestFetchForecastDataHTTPStatusError() {
	city := models.City{
		Id:        1,
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s", city.Latitude, city.Longitude, suite.apiKey),
		httpmock.NewStringResponder(500, "Internal Server Error"))

	result, err := suite.service.FetchForecastData(city, suite.apiKey)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func TestForecastServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ForecastServiceTestSuite))
}
