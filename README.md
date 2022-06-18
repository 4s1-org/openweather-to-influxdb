# OpenWeather to InfluxDb

## Build

```bash
make
```

## Configuration

Create a configuration file and save it to `config.json`.

```json
{
  "influxDb": {
    "serverUrl": "https://...:8086",
    "token": "...",
    "bucket": "...",
    "org": "...",
    "measurement": "..."
  }
}
```

## Start

```bash
find ../openweather-data/frankfurt/2022/06/* -type f -exec bin/openweather-to-influxdb -c config.json -s Frankfurt {} +
```
