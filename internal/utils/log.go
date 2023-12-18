package utils

import (
	"bytes"
	"fmt"
	logger "github.com/sirupsen/logrus"
)

var (
	Logger = logger.New()
)

func init() {
	Logger.SetFormatter(&ftr{})
}

type ftr struct {
}

func (f *ftr) Format(entry *logger.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("01-02 15:04:05")
	var newLog string
	newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}
