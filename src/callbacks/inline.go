package callbacks

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/SV1Stail/weather_bot/apiconnect"
	"github.com/SV1Stail/weather_bot/weather"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var User = struct {
	rmu       sync.RWMutex
	State     map[int64]bool
	WeatherID int
}{
	State: make(map[int64]bool),
}

func Callback(update *tgbotapi.Update) error {
	bot := apiconnect.GetBot()

	if update.CallbackQuery != nil {
		// Обрабатываем нажатие на инлайн-кнопку
		callbackData := update.CallbackQuery.Data

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введите координаты в формате:\nширота долгота")
		if _, err := bot.Send(msg); err != nil {
			slog.Error("send msg FAIL", "error", err)
			return err
		}
		if callbackData == "weather" {
			User.rmu.Lock()
			User.State[update.CallbackQuery.From.ID] = true
			User.WeatherID = 3
			User.rmu.Unlock()
		} else if callbackData == "weatherALL" {
			User.rmu.Lock()
			User.State[update.CallbackQuery.From.ID] = true
			User.WeatherID = 1
			User.rmu.Unlock()
		} else if callbackData == "weatherHours" {
			User.rmu.Lock()
			User.State[update.CallbackQuery.From.ID] = true
			User.WeatherID = 2
			User.rmu.Unlock()
		}
	}
	return nil
}

func GiveWeatherByCoordinats(update *tgbotapi.Update) {

	slog.Info("Messages start")
	bot := apiconnect.GetBot()
	lat, lon, err := askLatLon(GetNextMessage(update))
	slog.Info("askLatLon done", "lat", lat, "lon", lon)
	if err != nil {
		slog.Error("askLatLon FAIL", "error", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат координат")

		addWeatherKeyboard(&msg)
		if _, err := bot.Send(msg); err != nil {
			slog.Error("send error msg FAIL", "error", err)
			return
		}
		return
	}
	slog.Info("weatehr api request START")
	var weather weather.WeatherResponse
	weather.Get(lat, lon)
	slog.Info("weatehr api request SUCCESS")
	var weatherStr string
	if User.WeatherID == 1 {
		weatherStr = constructWeatherALL(&weather)
	} else if User.WeatherID == 2 {
		weatherStr = constructWeatherHours(&weather)
	} else if User.WeatherID == 3 {
		weatherStr = constructWeather(&weather)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, weatherStr)
	addWeatherKeyboard(&msg)
	if _, err := bot.Send(msg); err != nil {
		slog.Error("send WEATHER msg FAIL", "error", err)
		return
	}

}

func askLatLon(str string) (float64, float64, error) {
	slog.Info("input data", "str", str)
	coordinates := strings.Split(str, " ")
	slog.Info("check what in msg", "coordinates", coordinates, "str", str)
	if len(coordinates) < 2 || len(coordinates) > 2 {
		slog.Error("invalid number of coordinates", "coordinates", coordinates)
		return 0, 0, fmt.Errorf("need 2 numbers")
	}
	lat, err := strconv.ParseFloat(coordinates[0], 64)
	slog.Info("check lat", "lat", lat)
	if err != nil {
		slog.Error("lat convert FAIL", "error", err)
		return 0, 0, err
	}
	lon, err := strconv.ParseFloat(coordinates[1], 64)
	slog.Info("check lon", "lat", lon)
	if err != nil {
		slog.Error("lat convert FAIL", "error", err)
		return 0, 0, err
	}
	return lat, lon, nil
}
func constructWeatherALL(weather *weather.WeatherResponse) string {
	weatherStr := fmt.Sprintf(
		`Дата: %s
		Температура: %d°C
		Ощущается как: %d°C
		Условия: %s
		Скорость ветра: %.2f м/с`, weather.NowDt,
		weather.Fact.Temp,
		weather.Fact.FeelsLike,
		weather.Fact.Condition,
		weather.Fact.WindSpeed,
	)

	for _, f := range weather.Forecasts {
		for i, h := range f.Hours {
			if i < 10 {
				weatherStr += fmt.Sprintf(`
				Время 0%s:00:
				  температура: %d°C, 
				  ощущается:%d°C, 
				  ветер: %.2f м/с, 
				  состояние: %s
				`, h.Hour, h.Temp, h.FeelsLike, h.WindSpeed, h.Condition)
			} else {
				weatherStr += fmt.Sprintf(`
				Время %s:00: 
				  температура: %d°C, 
				  ощущается:%d°C, 
				  ветер: %.2f м/с, 
				  состояние: %s
				`, h.Hour, h.Temp, h.FeelsLike, h.WindSpeed, h.Condition)
			}
		}
	}
	weatherStr += "\n--------------------\n"
	for _, f := range weather.Forecasts {
		for k, v := range f.Parts {
			weatherStr += fmt.Sprintf(`
			Часть суток %s:\n
			   температура: %d°C, 
			   ощущается:%d°C, 
			   ветер: %.2f м/с, 
			   состояние: %s
			`, k, v.Temp, v.FeelsLike, v.WindSpeed, v.Condition)

		}
	}
	return weatherStr
}
func constructWeatherHours(weather *weather.WeatherResponse) string {
	weatherStr := ""
	for _, f := range weather.Forecasts {
		for i, h := range f.Hours {
			if i < 10 {
				weatherStr += fmt.Sprintf(`
				Время суток 0%s:00:
				   температура: %d°C, 
				   ощущается:%d°C, 
				   ветер: %.2f м/с, 
				   состояние: %s
				`, h.Hour, h.Temp, h.FeelsLike, h.WindSpeed, h.Condition)
			} else {
				weatherStr += fmt.Sprintf(`
				Время суток %s:00: 
				   температура: %d°C, 
				   ощущается:%d°C, 
				   ветер: %.2f м/с, 
				   состояние: %s
				`, h.Hour, h.Temp, h.FeelsLike, h.WindSpeed, h.Condition)
			}
		}
	}
	return weatherStr
}
func constructWeather(weather *weather.WeatherResponse) string {
	return fmt.Sprintf(
		`Дата: %s
	Температура: %d°C
	Ощущается как: %d°C
	Условия: %s
	Скорость ветра: %.2f м/с`, weather.NowDt,
		weather.Fact.Temp,
		weather.Fact.FeelsLike,
		weather.Fact.Condition,
		weather.Fact.WindSpeed,
	)
}

// capture next msg
func GetNextMessage(update *tgbotapi.Update) string {
	userID := update.Message.From.ID
	User.rmu.RLock()
	state, ok := User.State[userID]
	User.rmu.RUnlock()
	if ok && state {
		User.rmu.Lock()
		delete(User.State, userID)
		User.rmu.Unlock()
		return update.Message.Text
	}
	return ""
}

func addWeatherKeyboard(msg *tgbotapi.MessageConfig) {
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вся погода", "weatherALL"),        // 1
			tgbotapi.NewInlineKeyboardButtonData("Погода по часам", "weatherHours"), // 2
			tgbotapi.NewInlineKeyboardButtonData("Краткая погода", "weather"),       // 3
		),
	)
	msg.ReplyMarkup = keyboard
}
