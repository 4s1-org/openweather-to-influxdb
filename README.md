# OpenWeather to InfluxDB

## Build

```bash
make
```

## Configuration

Create a configuration file and save it to `config.json`.

```json
{
  "influxDB": {
    "serverUrl": "https://...:8086",
    "token": "...",
    "bucket": "...",
    "org": "...",
    "measurement": "..."
  }
}
```

## Start

The `-s` parameters is used for the city tag in InfluxDB

```bash
find ../openweather-data/frankfurt/2022/06/* -type f -exec bin/openweather-to-influxdb -c config.json -s Frankfurt {} +
```
