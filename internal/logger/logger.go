package logger

import (
	"log"
	"os"
)

type Logger struct {
	L  *log.Logger
	Nf *os.File
}

func New(logpath string) (logger *Logger, err error) {
	nf, err := os.Create(logpath)
	if err != nil {
		return
	}
	logger = &Logger{
		L: log.New(nf, "LOG", 19),
		Nf: nf,
	}
	
	return
}
