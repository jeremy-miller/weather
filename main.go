package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func main() {
	http.HandleFunc("/weather/", getWeather)
	http.ListenAndServe(":8080", nil)
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	data, error := query(city)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func query(city string) (weatherData, error) {
	APIkey := "c97e896907f938993a8b424cc8e026e0"
	res, error := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + APIkey + "&q=" + city)
	if error != nil {
		return weatherData{}, error
	}

	defer res.Body.Close()

	var data weatherData

	if error := json.NewDecoder(res.Body).Decode(&data); error != nil {
		return weatherData{}, error
	}

	return data, nil
}
