package cityservice

import (
	"encoding/json"
	"fmt"
	"testing"
	"weather-app/internal/models"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockCityRepository struct {
	mock.Mock
}

func (m *MockCityRepository) CreateCity(city models.City) (int, error) {
	args := m.Called(city)
	return args.Int(0), args.Error(1)
}

func (m *MockCityRepository) GetCities() ([]models.City, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.City), args.Error(1)
}

func (m *MockCityRepository) GetCity(cityId int) (models.City, error) {
	args := m.Called(cityId)
	return args.Get(0).(models.City), args.Error(1)
}

type CityServiceTestSuite struct {
	suite.Suite
	service  *CityService
	mockRepo *MockCityRepository
	apiKey   string
}

func (suite *CityServiceTestSuite) SetupTest() {
	suite.mockRepo = new(MockCityRepository)
	suite.service = NewCityService(suite.mockRepo)
	suite.apiKey = "test-api-key"
	httpmock.Activate()
}

func (suite *CityServiceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (suite *CityServiceTestSuite) TestCreateCity() {
	city := models.City{
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	suite.mockRepo.On("CreateCity", city).Return(1, nil)

	id, err := suite.service.CreateCity(city)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, id)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *CityServiceTestSuite) TestCreateCityError() {
	city := models.City{
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	suite.mockRepo.On("CreateCity", city).Return(0, fmt.Errorf("insert error"))

	id, err := suite.service.CreateCity(city)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, id)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *CityServiceTestSuite) TestGetCities() {
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

	suite.mockRepo.On("GetCities").Return(cities, nil)

	result, err := suite.service.GetCities()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), cities, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *CityServiceTestSuite) TestGetCitiesError() {
	suite.mockRepo.On("GetCities").Return(nil, fmt.Errorf("select error"))

	result, err := suite.service.GetCities()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *CityServiceTestSuite) TestGetCity() {
	city := models.City{
		Id:        1,
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	suite.mockRepo.On("GetCity", 1).Return(city, nil)

	result, err := suite.service.GetCity(1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), city, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *CityServiceTestSuite) TestGetCityError() {
	suite.mockRepo.On("GetCity", 1).Return(models.City{}, fmt.Errorf("select error"))

	result, err := suite.service.GetCity(1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), models.City{}, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *CityServiceTestSuite) TestFetchCityData() {
	cityName := "London"
	geoResponse := []geocodingResponse{
		{
			Name:    "London",
			Country: "GB",
			Lat:     51.5074,
			Lon:     -0.1278,
		},
	}
	responseBody, _ := json.Marshal(geoResponse)
	httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, suite.apiKey),
		httpmock.NewStringResponder(200, string(responseBody)))

	expectedCity := models.City{
		Name:      "London",
		Country:   "GB",
		Latitude:  51.5074,
		Longitude: -0.1278,
	}

	result, err := suite.service.FetchCityData(cityName, suite.apiKey)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedCity, result)
}

func (suite *CityServiceTestSuite) TestFetchCityDataNoResults() {
	cityName := "UnknownCity"
	httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, suite.apiKey),
		httpmock.NewStringResponder(200, "[]"))

	result, err := suite.service.FetchCityData(cityName, suite.apiKey)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), models.City{}, result)
	assert.Equal(suite.T(), fmt.Errorf("no results found for city: %s", cityName), err)
}

func (suite *CityServiceTestSuite) TestFetchCityDataError() {
	cityName := "London"
	httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, suite.apiKey),
		httpmock.NewErrorResponder(fmt.Errorf("network error")))

	result, err := suite.service.FetchCityData(cityName, suite.apiKey)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), models.City{}, result)
}

func (suite *CityServiceTestSuite) TestFetchCityDataInvalidJSON() {
	cityName := "London"
	httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, suite.apiKey),
		httpmock.NewStringResponder(200, "{invalid json}"))

	result, err := suite.service.FetchCityData(cityName, suite.apiKey)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), models.City{}, result)
}

func (suite *CityServiceTestSuite) TestFetchCityDataHTTPStatusError() {
	cityName := "London"
	httpmock.RegisterResponder("GET", fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, suite.apiKey),
		httpmock.NewStringResponder(500, "Internal Server Error"))

	result, err := suite.service.FetchCityData(cityName, suite.apiKey)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), models.City{}, result)
}

func TestCityServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CityServiceTestSuite))
}
