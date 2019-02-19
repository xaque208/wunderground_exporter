package exporter

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/xaque208/gowu"
)

type Forecast struct {
	High float32
	Low  float32
}

var (
	forecastHighTemp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wunderground_forecast_high_temperature",
		Help: "Temperature in Celcius",
	}, []string{"day"})

	forecastLowTemp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wunderground_forecast_low_temperature",
		Help: "Temperature in Celcius",
	}, []string{"day"})

	moonRiseTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wunderground_moonrise_time",
		Help: "Time of Moon Rise",
	}, nil)

	moonSetTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wunderground_moonrise_set",
		Help: "Time of Moon Set",
	}, nil)

	sunRiseTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wunderground_sunrise_time",
		Help: "Time of Sun Rise",
	}, nil)

	sunSetTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wunderground_sunrise_set",
		Help: "Time of Sun Set",
	}, nil)
)

func init() {
	prometheus.MustRegister(
		forecastHighTemp,
		forecastLowTemp,
		moonRiseTime,
		moonSetTime,
		sunRiseTime,
		sunSetTime,
	)
}

func forecastWatch(apiKey string) {
	c := gowu.NewClient(apiKey)
	fore, err := c.GetForecast("portland", "or")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, day := range fore.Simpleforecast.Forecastday {
		dayString := strconv.Itoa(i)

		highTemp, err := strconv.ParseFloat(day.High.Celsius, 32)
		if err != nil {
			log.Error(err)
		}

		lowTemp, err := strconv.ParseFloat(day.Low.Celsius, 32)
		if err != nil {
			log.Error(err)
		}

		forecastHighTemp.With(prometheus.Labels{"day": dayString}).Set(highTemp)
		forecastLowTemp.With(prometheus.Labels{"day": dayString}).Set(lowTemp)
	}
}

func astroWatch(apiKey string) {
	c := gowu.NewClient(apiKey)
	moonPhase, sunPhase, err := c.GetAstronomy("portland", "or")
	if err != nil {
		fmt.Println(err)
		return
	}

	moonRiseHourMin, err := strconv.ParseFloat(
		fmt.Sprintf("%s.%s", moonPhase.MoonRise.Hour, moonPhase.MoonRise.Minute), 32)
	moonSetHourMin, err := strconv.ParseFloat(
		fmt.Sprintf("%s.%s", moonPhase.MoonSet.Hour, moonPhase.MoonSet.Minute), 32)

	sunRiseHourMin, err := strconv.ParseFloat(
		fmt.Sprintf("%s.%s", sunPhase.SunRise.Hour, sunPhase.SunRise.Minute), 32)
	sunSetHourMin, err := strconv.ParseFloat(
		fmt.Sprintf("%s.%s", sunPhase.SunSet.Hour, sunPhase.SunSet.Minute), 32)

	moonRiseTime.With(prometheus.Labels{}).Set(moonRiseHourMin)
	moonSetTime.With(prometheus.Labels{}).Set(moonSetHourMin)

	sunRiseTime.With(prometheus.Labels{}).Set(sunRiseHourMin)
	sunSetTime.With(prometheus.Labels{}).Set(sunSetHourMin)
}

func ScrapeMetrics(apiKey string) {
	forecastWatch(apiKey)
	astroWatch(apiKey)
}
