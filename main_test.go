package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWeather(t *testing.T) {
	req, err := http.NewRequest("GET", "/weather/tokyo", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getWeather)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusOK)
	}
	var res struct {
		City string  `json:"city"`
		Temp float64 `json:"temp"`
		Took string  `json:"took"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Error(err)
	}
	if res.City != "tokyo" {
		t.Errorf("handler returned incorrect city: %s", res.City)
	}
}
