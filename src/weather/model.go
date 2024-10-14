package weather

var keyWeather = ""
//b5d33b80-8441-4e9a-a6fb-2118d1d9e1b5
type WeatherResponse struct {
	Now       int64      `json:"now"`
	NowDt     string     `json:"now_dt"`
	Info      Info       `json:"info"`
	Fact      Fact       `json:"fact"`
	Forecasts []Forecast `json:"forecasts"`
}
type Info struct {
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	URL    string  `json:"url"`
	Tzinfo Tzinfo  `json:"tzinfo"`
}
type Tzinfo struct {
	Name string `json:"name"`
}

type Forecast struct {
	Date  string          `json:"date"`
	Parts map[string]Fact `json:"parts"` // key = часть суток val = data
	Hours []Fact          `json:"hours"` // key = час суток val = data

}
type Fact struct {
	Hour      string  `json:"hour,omitempty"` // only for hours
	Temp      int     `json:"temp"`
	FeelsLike int     `json:"feels_like"`
	Condition string  `json:"condition"`
	WindSpeed float32 `json:"wind_speed"`
}
