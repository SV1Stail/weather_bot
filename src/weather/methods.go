package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func (w *WeatherResponse) Get(lat, lon float64) error {
	url := fmt.Sprintf("https://api.weather.yandex.ru/v2/forecast?lat=%f&lon=%f&lang=ru_RU&limit=1", lat, lon)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("cant build new request", "error", err)
		return err
	}
	req.Header.Add("X-Yandex-API-Key", keyWeather)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("request FAIL", "error", err)
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("read body FAIL", "error", err)
		return err
	}

	if err := json.Unmarshal(body, w); err != nil {
		slog.Error("unmarshal FAIL", "error", err)
		return err
	}

	for _, f := range w.Forecasts {
		for k, v := range f.Parts {
			fmt.Printf("[%s]: temp: %d, feel: %d, cond: %s\n", k, v.Temp, v.FeelsLike, v.Condition)
		}
	}
	for _, s := range w.Forecasts {
		for _, f := range s.Hours {
			fmt.Printf("[%s]: temp: %d, feel: %d, cond: %s\n", f.Hour, f.Temp, f.FeelsLike, f.Condition)
		}
	}
	return nil
}
