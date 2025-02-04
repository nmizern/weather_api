package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/nmizern/weather-api/internal/services"
)


type WeatherHandler struct {
	WeatherService services.WeatherService
	Logger         *zap.Logger
}


func NewWeatherHandler(ws services.WeatherService, logger *zap.Logger) *WeatherHandler {
	return &WeatherHandler{
		WeatherService: ws,
		Logger:         logger,
	}
}


func (h *WeatherHandler) GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city query parameter is required"})
		return
	}

	weather, err := h.WeatherService.GetWeather(city)
	if err != nil {
		h.Logger.Error("Failed to get weather", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve weather data"})
		return
	}

	c.JSON(http.StatusOK, weather)
}
