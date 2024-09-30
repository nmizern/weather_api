package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "Choose the city", http.StatusBadRequest)
		return
	}

	apiKey := "88d0d1105d17ddab6da71f6245005636"
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to load weather data", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to load weather data", resp.StatusCode)
		return
	}

	var weatherData struct {
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity int     `json:"humidity"`
		} `json:"main"`
		Name string `json:"name"`
	}

	err = json.NewDecoder(resp.Body).Decode(&weatherData)
	if err != nil {
		http.Error(w, "Problem with parsing", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"city":        weatherData.Name,
		"temperature": weatherData.Main.Temp,
		"humidity":    weatherData.Main.Humidity,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/weather", weatherHandler)
	log.Println("Server is starting")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
