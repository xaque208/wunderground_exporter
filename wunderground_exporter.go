package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xaque208/gowu"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
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
)

func forecastWatch(apiKey string) {

	for {
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

		//Sleep 15 minutes between updates for API limits
		time.Sleep(900 * time.Second)
	}
}

func init() {
	prometheus.MustRegister(forecastHighTemp)
	prometheus.MustRegister(forecastLowTemp)
}

func main() {

	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9101").String()
		configPath    = kingpin.Flag("config", "Specify the configuration directory").Default(".").String()
		verbose       = kingpin.Flag("verbose", "Increase verbosity").Short('v').Bool()
	)

	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	viper.SetConfigName("wunderground_exporter")
	viper.AddConfigPath(*configPath)
	viper.AddConfigPath(".")

	log.Debug("Reading configuration file")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	wuApiKey := viper.GetString("wunderground.apikey")

	go forecastWatch(wuApiKey)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
