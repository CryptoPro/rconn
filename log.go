package main

import (
	"os"

	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("rconn")

func initLogger() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendError := logging.NewLogBackend(os.Stderr, "", 0)
	backendErrorLevered := logging.AddModuleLevel(backendError)
	backendErrorLevered.SetLevel(logging.ERROR, "")

	backendFormatter := logging.NewBackendFormatter(backend, logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	))
	if false {
		logging.SetBackend(backendFormatter, backendErrorLevered)
	} else {
		logging.SetBackend(backendErrorLevered)
	}

}
