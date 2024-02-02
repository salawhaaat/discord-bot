package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WeatherResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Country        string  `json:"country"`
		Localtime      string  `json:"localtime"`
		TzID           string  `json:"tz_id"`
		Latitude       float64 `json:"lat"`
		Longitude      float64 `json:"lon"`
		LocaltimeEpoch int64   `json:"localtime_epoch"`
	} `json:"location"`

	Current struct {
		LastUpdatedEpoch int64   `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TemperatureC     float64 `json:"temp_c"`
		Condition        struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

func GetWeather(apiKey, city string) (string, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city)

	// Make HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "Invalid city", err
	}

	// Decode JSON response
	var weatherResponse WeatherResponse
	err = json.NewDecoder(response.Body).Decode(&weatherResponse)
	if err != nil {
		return "", err
	}

	// Ensure Localtime has enough characters before extracting date and time
	localTime := weatherResponse.Location.Localtime
	date, time := "", ""
	if len(localTime) >= 16 {
		date = localTime[:10]
		time = localTime[11:16]
	}

	result := struct {
		City      string
		Temp      float64
		Date      string
		Time      string
		Condition string
	}{
		City:      weatherResponse.Location.Name,
		Temp:      weatherResponse.Current.TemperatureC,
		Date:      date,
		Time:      time,
		Condition: weatherResponse.Current.Condition.Text,
	}

	resultString := fmt.Sprintf("🗺 City: %s\n🌡 Temperature: %.1f°C\n🪟 Condition: %s\n🗓 Date: %s\n🕰 Time: %s",
		result.City, result.Temp, result.Condition, result.Date, result.Time)

	return resultString, nil
}
