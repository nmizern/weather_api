package main

import (
	"log"
	
	"github.com/gin-gonic/gin"
	"github.com/nmizern/weather-api/internal/config"
	"github.com/nmizern/weather-api/internal/handlers"
	"github.com/nmizern/weather-api/internal/kafka"
	"github.com/nmizern/weather-api/internal/services"
	"github.com/nmizern/weather-api/internal/utils"
	"go.uber.org/zap"
)

func main() {
	
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	
	logger, err := utils.NewLogger()
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logger.Sync()

	
	db, err := config.ConnectPostgres(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}
	
	db.AutoMigrate(&services.WeatherRecord{})

	
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers, "weather_events")
	defer kafkaProducer.Close()

	
	weatherService := services.NewWeatherService(cfg, db, kafkaProducer, logger)

	
	router := gin.Default()
	weatherHandler := handlers.NewWeatherHandler(weatherService, logger)
	router.GET("/weather", weatherHandler.GetWeather)

	
	logger.Info("Starting server", utils.String("port", cfg.Port))
	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Fatal("Error starting server", zap.Error(err))
	}
}
