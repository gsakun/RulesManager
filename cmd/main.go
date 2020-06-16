package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gsakun/RulesManager/db"
	"github.com/gsakun/RulesManager/handler"
	colorable "github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	loglevel := os.Getenv("LOG_LEVEL")
	var logLevel log.Level
	log.Infof("loglevel env is %s", loglevel)
	if loglevel == "debug" {
		log.SetLevel(log.DebugLevel)
		logLevel = log.DebugLevel
		log.Infof("log level is %s", loglevel)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
		logLevel = log.InfoLevel
		log.Infoln("log level is normal")
	}
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/rulemanager.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     5, //days
		Level:      logLevel,
		Formatter: &log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	})
	log.SetOutput(colorable.NewColorableStdout())
	if err != nil {
		log.Fatalf("Failed to initialize file rotate hook: %v", err)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
	log.SetReportCaller(true)
	log.AddHook(rotateFileHook)
}

func main() {
	var (
		dbaddress = kingpin.Flag(
			"database",
			"The database address for get machine info.",
		).Default("").String()
		maxconn = kingpin.Flag(
			"maxconn",
			"Database maxconn.",
		).Default("100").Int()
		maxidle = kingpin.Flag(
			"maxidle",
			"Database maxidle.",
		).Default(string(*maxconn)).Int()
		requiredgroup = kingpin.Flag("requiredgroup", "this region required type. for example 1,2,3").Default("").String()
		rulespath     = kingpin.Flag("rulespath", "Rule store path").Default("").String()
		interval      = kingpin.Flag("interval", "Sync Interval.").Default("60").Int64()
	)

	kingpin.Version("RulerManager v1.0")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	go func() {
		db, err := db.Init(*dbaddress, *maxconn, *maxidle)
		if err != nil {
			log.Errorf("ping db fail:%v", err)
			time.Sleep(60 * time.Second)
		} else {
			defer db.Close()
			//log.Infoln("START SYNC")
			handler.Queryalertgroup(db)
			var path string
			if strings.HasSuffix(*rulespath, "/") {
				path = *rulespath
			} else {
				path = fmt.Sprintf("%s/", *rulespath)
			}
			err := handler.HandlerRule(db, *requiredgroup, path)
			if err != nil {
				time.Sleep(60 * time.Second)
			}
			time.Sleep(time.Duration(*interval) * time.Second)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		os.Exit(0)
	}()
	select {}
}
