# influx-loader
Simple utility for loading data into influx/telegraf

## How to Install
`go get github.com/abferm/influx-loader`

## Usage
```
$ influx-loader -h
NAME:
   influx-loader - Uploads various log files to InfluxDB via telegraf

USAGE:
   influx-loader [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   mooshimeter, mm  import mooshimeter log
   metabase, mb     import metabase log
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value  Path to log file (default: "./log.csv")
   --telegraf value        telegraf address [host:port] (default: "localhost:8094")
   --tag value, -t value   Tag [key:value]
   --debug                 Print metrics to screen rather than sending to telegraf
   --help, -h              show help
   --version, -v           print the version
```
