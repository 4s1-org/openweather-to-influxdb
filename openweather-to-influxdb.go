package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type ConfigurationInfluxDB struct {
	ServerUrl   string `json:"serverUrl"`
	Token       string `json:"token"`
	Bucket      string `json:"bucket"`
	Org         string `json:"org"`
	Measurement string `json:"measurement"`
}

type Configuration struct {
	InfluxDB ConfigurationInfluxDB `json:"influxDB"`
}

type OpenWeatherDataCoord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type OpenWeatherDataWeather struct {
	Id          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type OpenWeatherDataMain struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int64   `json:"pressure"`
	Humidity  int64   `json:"humidity"`
}

type OpenWeatherDataSys struct {
	Type    int64  `json:"type"`
	Id      int64  `json:"id"`
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
}

type OpenWeatherDataClouds struct {
	All int64 `json:"all"`
}

type OpenWeatherDataWind struct {
	Speed float64 `json:"speed"`
	Deg   int64   `json:"deg"`
	Gust  float64 `json:"gust"`
}

type OpenWeatherData struct {
	Coord      OpenWeatherDataCoord     `json:"coord"`
	Weather    []OpenWeatherDataWeather `json:"weather"`
	Base       string                   `json:"base"`
	Main       OpenWeatherDataMain      `json:"main"`
	Visibility int64                    `json:"visibility"`
	Wind       OpenWeatherDataWind      `json:"wind"`
	Clouds     OpenWeatherDataClouds    `json:"clouds"`
	Dt         int64                    `json:"dt"`
	Sys        OpenWeatherDataSys       `json:"sys"`
	Timezone   int64                    `json:"timezone"`
	Id         int64                    `json:"id"`
	Name       string                   `json:"name"`
	Cod        int64                    `json:"cod"`
}

func main() {
	// Program parameters
	var configFileName string
	flag.StringVar(&configFileName, "c", "", "configuration file")

	var city string
	flag.StringVar(&city, "s", "", "city")

	flag.Parse()

	sourceFiles := flag.Args()

	// Load configuration
	config, err := loadConfigurationFile(configFileName)
	if err != nil {
		panic(err)
	}

	// create new client with default option for server url authenticate by token
	client := influxdb2.NewClientWithOptions(
		config.InfluxDB.ServerUrl,
		config.InfluxDB.Token,
		influxdb2.DefaultOptions().SetBatchSize(20))
	// user blocking write client for writes to desired bucket
	writeAPI := client.WriteAPI(config.InfluxDB.Org, config.InfluxDB.Bucket)

	for _, filename := range sourceFiles {
		fmt.Println(filename)
		doIt(config, filename, writeAPI, city)
	}

	writeAPI.Flush()
	client.Close()
}

func loadConfigurationFile(configFileName string) (*Configuration, error) {
	_, err := os.Stat(configFileName)
	if errors.Is(err, os.ErrNotExist) {
		// Config doesn't exists
		return nil, errors.New("Configuration file not found")
	}
	if err != nil {
		panic(err)
	}

	file, err := ioutil.ReadFile(configFileName)
	if err != nil {
		panic(err)
	}
	configuration := Configuration{}
	err = json.Unmarshal([]byte(file), &configuration)
	if err != nil {
		panic(err)
	}

	return &configuration, nil
}

func doIt(config *Configuration, filename string, writeAPI api.WriteAPI, city string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	openWeatherData := OpenWeatherData{}
	err = json.Unmarshal([]byte(file), &openWeatherData)
	if err != nil {
		panic(err)
	}

	point := influxdb2.NewPointWithMeasurement(config.InfluxDB.Measurement)
	point.
		AddTag("Stadt", city).
		AddField("Temperatur", openWeatherData.Main.Temp).
		AddField("Temperatur (gefühlt)", openWeatherData.Main.FeelsLike).
		AddField("Luftdruck", openWeatherData.Main.Pressure).
		AddField("Windgeschwindigkeit", openWeatherData.Wind.Speed).
		SetTime(time.Unix(openWeatherData.Dt, 0))
	writeAPI.WritePoint(point)
}
