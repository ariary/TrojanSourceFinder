package utils

import (
	"log"
	"os"
)

//LOGGER
var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func InitLoggers() {
	InfoLogger = log.New(os.Stdout, "", 0)
	ErrorLogger = log.New(os.Stderr, "", 0)
}
