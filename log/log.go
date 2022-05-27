package log

import (
	"log"
	"os"
)

/*
func InitLogger() (l *log.Logger, err error) {

	nf, err := os.Create("logs")

	if err != nil {
		fmt.Println("loger error")
		return
	}

	l = log.New(nf, "LOG", 19)

	l.Printf("hello")

	return

}

*/

func InitLogger() (logger *log.Logger) {
		
	logger = log.New(os.Stdout, "LOG", 19)

	return logger
}

/*

func InitLogger() (logger *log.Logger, err error) {

	nf, err := os.Create("logs")

	if err != nil {
		fmt.Println("loger error")
		return
	}
		
	Logger = log.New(nf, "LOG", 19)

	Logger.Printf("hello")

	return

}

*/