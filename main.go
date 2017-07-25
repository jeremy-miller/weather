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
	temp, err := mw.temperature(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
			temp, err := p.temperature(city)
			if err != nil {
				errors <- err
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
		case err := <-errors:
			return 0, err
		}
	}
	return sum / float64(len(w)), nil
}

type openWeatherMap struct {
	apiKey string
}

func (w openWeatherMap) temperature(city string) (float64, error) {
	res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + w.apiKey + "&q=" + city)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	var data struct { // deconstruct the response and only pull out the temperature
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return 0, err
	}
	log.Printf("openWeathermap: %s = %.2fK", city, data.Main.Kelvin)
	return data.Main.Kelvin, nil
}
