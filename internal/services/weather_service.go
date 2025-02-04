package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/nmizern/weather-api/internal/config"
	"github.com/nmizern/weather-api/internal/models"
	"github.com/nmizern/weather-api/internal/kafka"
)


type WeatherRecord = models.WeatherRecord


type WeatherService interface {
	GetWeather(city string) (*WeatherRecord, error)
}

type weatherServiceImpl struct {
	Config     *config.Config
	DB         *gorm.DB
	Producer   *kafka.Producer
	Logger     *zap.Logger
	HTTPClient *http.Client
}


func NewWeatherService(cfg *config.Config, db *gorm.DB, producer *kafka.Producer, logger *zap.Logger) WeatherService {
	return &weatherServiceImpl{
		Config:     cfg,
		DB:         db,
		Producer:   producer,
		Logger:     logger,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}


type OpenWeatherResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Name string `json:"name"`
}

func (s *weatherServiceImpl) GetWeather(city string) (*WeatherRecord, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, s.Config.APIKey)
	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call weather API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status: %d", resp.StatusCode)
	}

	var weatherData OpenWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("failed to decode weather API response: %w", err)
	}

	record := &WeatherRecord{
		City:        weatherData.Name,
		Temperature: weatherData.Main.Temp,
		Humidity:    weatherData.Main.Humidity,
		CreatedAt:   time.Now(),
	}

	
	if err := s.DB.Create(record).Error; err != nil {
		s.Logger.Error("Failed to save weather record", zap.Error(err))
		
	}

	
	msg := fmt.Sprintf("Weather data for %s recorded at %s", record.City, record.CreatedAt.Format(time.RFC3339))
	if err := s.Producer.SendMessage("weather_events", msg); err != nil {
		s.Logger.Error("Failed to send Kafka message", zap.Error(err))
		
	}

	return record, nil
}
