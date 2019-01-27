package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

func EnsureDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func Log(msg ...interface{}) {
	log.Output(2, fmt.Sprintln(msg...))
}

func Logf(format string, msg ...interface{}) {
	log.Output(2, fmt.Sprintf(format, msg...))
}

func GetLogWriter() io.Writer {
	return os.Stdout
}

func CloseLogWriter() {

}
