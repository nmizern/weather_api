package models

import "time"

type WeatherRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	City        string    `json:"city"`
	Temperature float64   `json:"temperature"`
	Humidity    int       `json:"humidity"`
	CreatedAt   time.Time `json:"created_at"`
}
