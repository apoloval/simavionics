package simavionics

import (
	"os"

	"github.com/op/go-logging"
)

const logFormat = `%{color}%{time:15:04:05.000} %{module:12s} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`

func EnableLogging() {
	EnableLoggingLevel(logging.INFO)
}

func EnableLoggingLevel(level logging.Level) {
	b1 := logging.NewLogBackend(os.Stdout, "", 0)
	formatter := logging.MustStringFormatter(logFormat)
	b2 := logging.NewBackendFormatter(b1, formatter)
	b3 := logging.AddModuleLevel(b2)
	b3.SetLevel(level, "")

	logging.SetBackend(b3)
}

func DisableLogging() {
	logging.SetBackend()
}
