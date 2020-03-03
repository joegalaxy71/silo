package helpers

import (
	"github.com/op/go-logging"
	"os"
	"os/signal"
	"syscall"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func HandleSignals(c chan os.Signal) {
	<-c
	os.Exit(1)
}

func InitSignals() {
	// handle ^c (os.Interrupt)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go HandleSignals(c)
}

func InitLogs(verbose bool) *logging.Logger {
	// logging
	// color logging for terminal logs
	formatColor := logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfunc} (%{shortfile}) ▶ %{level:.4s} %{id}%{color:reset} %{message}")
	// no colors for file/service logging
	formatBase := logging.MustStringFormatter(
		"%{time:15:04:05.000} %{shortfunc} (%{shortfile}) ▶ %{level:.4s} %{id} %{message}")

	// terminal backend write fo stderr
	backendTerminal := logging.NewLogBackend(os.Stderr, "", 0)
	backendTerminalFormatter := logging.NewBackendFormatter(backendTerminal, formatColor)
	backendTerminalLeveled := logging.AddModuleLevel(backendTerminalFormatter)
	if verbose == true {
		backendTerminalLeveled.SetLevel(logging.DEBUG, "")
	} else {
		backendTerminalLeveled.SetLevel(logging.INFO, "")
	}

	// file writes to logfile
	backendFile := logging.NewLogBackend(os.Stderr, "", 0)
	backendFileFormatter := logging.NewBackendFormatter(backendFile, formatBase)
	backendFileLeveled := logging.AddModuleLevel(backendFileFormatter)
	backendFileLeveled.SetLevel(logging.INFO, "")

	//logging.SetBackend(backendFileLeveled, backendTerminalLeveled)
	logging.SetBackend(backendTerminalLeveled)

	logger := logging.MustGetLogger("example")
	return logger
}
