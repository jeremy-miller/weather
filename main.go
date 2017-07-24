package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// List of all weather APIs to be queried, along with their API keys.
var mw = multiWeatherProvider{
	openWeatherMap{apiKey: "c97e896907f938993a8b424cc8e026e0"},
}

func main() {
	http.HandleFunc("/weather/", getWeather)
	http.ListenAndServe(":8080", nil)
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	temp, error := mw.temperature(city)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"city": city,
		"temp": temp,
		"took": time.Since(begin).String(),
	})
}

type weatherProvider interface {
	temperature(city string) (float64, error) // returns temperature in Kelvin
}

type multiWeatherProvider []weatherProvider

func (w multiWeatherProvider) temperature(city string) (float64, error) {
	// Add temperatures and errors to their own channels, to be processed later.
	temps := make(chan float64, len(w))
	errors := make(chan error, len(w))
	for _, provider := range w {
		go func(p weatherProvider) {
			temp, error := p.temperature(city)
			if error != nil {
				errors <- error
				return
			}
			temps <- temp
		}(provider)
	}
	sum := 0.0
	for i := 0; i < len(w); i++ {
		select {
		case temp := <-temps:
			sum += temp
		case error := <-errors:
			return 0, error
		}
	}
	return sum / float64(len(w)), nil
}

type openWeatherMap struct {
	apiKey string
}

func (w openWeatherMap) temperature(city string) (float64, error) {
	res, error := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + w.apiKey + "&q=" + city)
	if error != nil {
		return 0, error
	}
	defer res.Body.Close()
	var data struct { // deconstruct the response and only pull out the temperature
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}
	if error := json.NewDecoder(res.Body).Decode(&data); error != nil {
		return 0, error
	}
	log.Printf("openWeathermap: %s = %.2fK", city, data.Main.Kelvin)
	return data.Main.Kelvin, nil
}
