![Go workflow](https://github.com/tetragramato/go-conz/actions/workflows/go.yml/badge.svg)
[![GitHub release](https://img.shields.io/github/release/tetragramato/go-conz/all.svg)](https://github.com/Tetragramato/go-conz/releases)

# Go-Conz

Consume [Phoscon ZigBee API](https://dresden-elektronik.github.io/deconz-rest-doc/getting_started/) to:

- retrieve sensors information,
- store data in a DB (Badger),
- then serve via a HTTP Server.

For now, `go-conz` collect only [Xiaomi aqara](https://www.aqara.com/en/temperature_humidity_sensor.html) sensors data (
temperature, humidity and pressure), because i only own those !

# Installation and Run

You should download the right release for your system, and unzip the archive :

- Linux (arm, amd64, 386)
- Windows (386, amd64)
- Mac _darwin_ (amd64)

## Configuration

You can add a config file in `go-conz-config.yaml`, at the same level as `go-conz` app :

```yaml
# Keys with default values
phosconUrl: "https://phoscon.de/discover"
databasePath: "./goconz-sensors"
delayInSecond: 30
traceHttp: true
httpPort: ":9000"
readOnly: false
```

Or your can pass configuration values to the command line if you wish :

```
./go-conz --databasePath=./exemple_db
```

Or put environment variable (the `GOCONZ` prefix is mandatory) :

```
GOCONZ_PHOSCONURL=https://phoscon.de/discover
```

## Run

Just run with :

```
./go-conz
```

You can then open a browser and go to `http://localhost:9000/sensors` (if you have not changed the `httpPort`)

JSON structure example :

```json
[
  {
    "uniqueId": "00:15:8d:00:09:69:6c:da-01-0492",
    "name": "kitchen sensor",
    "type": "ZHATemperature",
    "events": [
      {
        "etag": "9963da7eaf696b7020095e5227a09335",
        "lastUpdated": "2021-06-21T09:36:07.058",
        "temperature": 2568,
        "humidity": 0,
        "pressure": 0
      },
      {
        "etag": "0049fe00bff6f9dc89b24020adf3f755",
        "lastUpdated": "2021-06-21T14:56:19.732",
        "temperature": 2637,
        "humidity": 0,
        "pressure": 0
      },
      {
        "etag": "5a11d1a8355846a36e75757ae9404fbb",
        "lastUpdated": "2021-06-21T15:03:51.792",
        "temperature": 2643,
        "humidity": 0,
        "pressure": 0
      }
    ]
  }
]
```

`events` are sorted by insertion order in database.

# Dependencies

- [Viper](https://github.com/spf13/viper) for configuration
- [Badger](https://github.com/dgraph-io/badger) for database
- [Resty](https://github.com/go-resty/resty) for REST http client & retry
