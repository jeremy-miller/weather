package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	mw := multiWeatherProvider{
		openWeatherMap{apiKey: "c97e896907f938993a8b424cc8e026e0"},
		weatherUnderground{apiKey: "d78b4951b88aac57"},
	}
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
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
	})
	http.ListenAndServe(":8080", nil)
}

type weatherProvider interface {
	temperature(city string) (float64, error) // in Kelvin
}

type multiWeatherProvider []weatherProvider

func (w multiWeatherProvider) temperature(city string) (float64, error) {
	sum := 0.0
	for _, provider := range w {
		temp, error := provider.temperature(city)
		if error != nil {
			return 0, error
		}
		sum += temp
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
	var data struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}
	if error := json.NewDecoder(res.Body).Decode(&data); error != nil {
		return 0, error
	}
	log.Printf("openWeathermap: %s: %.2f", city, data.Main.Kelvin)
	return data.Main.Kelvin, nil
}

type weatherUnderground struct {
	apiKey string
}

func (w weatherUnderground) temperature(city string) (float64, error) {
	res, error := http.Get("http://api.wunderground.com/api/" + w.apiKey + "/conditions/q/" + city + ".json")
	if error != nil {
		return 0, error
	}
	defer res.Body.Close()
	var data struct {
		Observation struct {
			Celcius float64 `json:"temp_c"`
		} `json:"current_observation"`
	}
	if error := json.NewDecoder(res.Body).Decode(&data); error != nil {
		return 0, error
	}
	kelvin := data.Observation.Celcius + 273.15
	log.Printf("weatherUnderground: %s: %.2f", city, kelvin)
	return kelvin, nil
}
