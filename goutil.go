package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"github.com/zhengcf/goutil/printer"
	"github.com/zhengcf/goutil/config"
	"github.com/zhengcf/goutil/util/logutil"
	"github.com/zhengcf/goutil/util/signal"
	"github.com/zhengcf/goutil/util/errors"
	log "github.com/sirupsen/logrus"
)


// Flag Names
const (
	nmVersion          = "V"
	nmHost             = "host"
	nmPort 			   = "P"
	nmPsAddr           = "psaddr"
	nmConfig           = "config"
	nmLogLevel         = "L"
	nmLogFile          = "log-file"
)

var (
	appReleaseVersion = ""
	version    = flagBoolean(nmVersion, false, "print version information and exit")
	configPath = flag.String(nmConfig, "", "config file path")

	host       = flag.String(nmHost, "0.0.0.0", "go-util server host")
	port       = flag.String(nmPort, "8000", "go-util server port")
	psaddr     = flag.String(nmPsAddr, "", "ps't addr for connect")
	// Log
	logLevel     = flag.String(nmLogLevel, "info", "log level: info, debug, warn, errors, fatal")
	logFile      = flag.String(nmLogFile, "", "log file path")

)

var (
	graceful bool

	cfg      *config.Config
	listener *net.Listener
)


func main() {

	flag.Parse()
	if *version {
		fmt.Println(printer.GetAppInfo())
		os.Exit(0)
	}

	loadConfig()
	overrideConfig()
	printInfo()

	setupLog()
	createServer()
	signal.SetupSignalHandler(serverShutdown)
	runSever()
	os.Exit(0)
}



func flagBoolean(name string, defaultVal bool, usage string) *bool {
	if defaultVal == false {
		// Fix #4125, golang do not print default false value in usage, so we append it.
		usage = fmt.Sprintf("%s (default false)", usage)
		return flag.Bool(name, defaultVal, usage)
	}
	return flag.Bool(name, defaultVal, usage)
}


func loadConfig() {
	cfg = config.GetGlobalConfig()
	if *configPath != "" {
		err := cfg.Load(*configPath)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func setupLog() {
	err := logutil.InitLogger(cfg.Log.ToLogConfig())
	if err != nil {
		log.Fatal(err.Error())
	}
}

func overrideConfig() {
	actualFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		actualFlags[f.Name] = true
	})

	// Base
	if actualFlags[nmHost] {
		cfg.Host = *host
	}
	var err error
	if actualFlags[nmPort] {
		var p int
		p, err = strconv.Atoi(*port)
		log.Fatal(err)
		cfg.Port = uint(p)
	}
	if actualFlags[nmPsAddr] {
		cfg.Psaddr = *psaddr
	}
	// Log
	if actualFlags[nmLogLevel] {
		cfg.Log.Level = *logLevel
	}
	if actualFlags[nmLogFile] {
		cfg.Log.File.Filename = *logFile
	}
}



func printInfo() {
	// Make sure the TiDB info is always printed.
	level := log.GetLevel()
	log.SetLevel(log.InfoLevel)
	printer.PrintAppInfo()
	log.SetLevel(level)
}


func serverShutdown(isgraceful bool) {

	if isgraceful {
		graceful = true
	}

	if listener != nil {
		lis := *listener
		lis.Close()
	}
}

func createServer() {
	addr := cfg.Host + ":" + strconv.FormatUint(uint64(cfg.Port), 10)
	_, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(errors.New("errors listen on :"+addr))
	}
}

func runSever() {
	if listener == nil {
		log.Fatal(errors.New("No availavle Listen"))
	}
}
