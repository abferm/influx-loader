package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"runtime/debug"
	"strings"

	"github.com/urfave/cli"
)

var log Log

func main() {
	app := cli.NewApp()
	app.Name = "influx-loader"
	app.Usage = "Uploads various log files to InfluxDB via telegraf"
	app.Commands = []cli.Command{mooshimeter, metabase}
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		app.Version = buildInfo.Main.Version
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file,f",
			Usage: "Path to log file",
			Value: "./log.csv",
		},
		cli.StringFlag{
			Name:  "telegraf",
			Usage: "telegraf address [host:port]",
			Value: "localhost:8094",
		},
		cli.StringSliceFlag{
			Name:  "tag,t",
			Usage: "Add tag to metrics, may be specified more than once to add multiple tags [key:value]",
		},
		cli.DurationFlag{
			Name:  "shift,s",
			Usage: "Shift points by duration",
			Value: 0,
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Print metrics to screen rather than sending to telegraf",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

var mooshimeter = cli.Command{
	Name:      "mooshimeter",
	Aliases:   []string{"mm"},
	Usage:     "import mooshimeter log",
	ArgsUsage: "[channel names...]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "metric,m",
			Value: "mooshimeter",
			Usage: "Influx metric name",
		},
	},
	Action: func(ctx *cli.Context) error {
		fileName := ctx.GlobalString("file")
		f, err := os.Open(fileName)
		if err != nil {
			return err
		}
		log = NewMooshimeterLog(f, ctx.Args())
		return nil
	},
	After: upload,
}

var metabase = cli.Command{
	Name:    "metabase",
	Aliases: []string{"mb"},
	Usage:   "import metabase log",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "metric,m",
			Value: "metabase",
			Usage: "Influx metric name",
		},
	},
	Action: func(ctx *cli.Context) error {
		fileName := ctx.GlobalString("file")
		f, err := os.Open(fileName)
		if err != nil {
			return err
		}
		log, err = NewMetaBaseLog(f)
		return err
	},
	After: upload,
}

var upload = func(ctx *cli.Context) error {
	metricName := ctx.String("metric")
	shift := ctx.GlobalDuration("shift")
	tags, err := toTags(ctx.GlobalStringSlice("tag"))
	if err != nil {
		return err
	}
	var conn io.WriteCloser
	conn, err = net.Dial("udp", ctx.GlobalString("telegraf"))
	if err != nil {
		return err
	}
	if ctx.GlobalBool("debug") {
		conn = os.Stdout
	}
	sent := 0
	defer func() {
		fmt.Printf("Successfully sent %d metrics\n", sent)
	}()
	for entry, err := log.Next(); ; entry, err = log.Next() {
		if err != nil {
			if err == io.EOF {
				return conn.Close()
			}
			return err
		}
		// shift entry time
		entry.Shift(shift)
		_, err := fmt.Fprintln(conn, entry.Telegraf(metricName, tags).ToLineProtocal())
		if err != nil {
			return err
		}
		sent++
	}
}

func toTags(tagEntries []string) (tags map[string]string, err error) {
	tags = make(map[string]string, len(tagEntries))
	for _, s := range tagEntries {
		split := strings.Split(s, ":")
		if len(split) != 2 {
			err = fmt.Errorf("invalid tag entry %q", s)
			return
		}
		tags[split[0]] = split[1]
	}
	return
}
