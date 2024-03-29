package logapp

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type ILogger interface {
	logrus.FieldLogger
}

type Config struct {
	Level         int
	Colors        bool
	FullTimestamp bool
}

func NewLogger(config *Config) ILogger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:     config.Colors,
		FullTimestamp:   config.FullTimestamp,
		TimestampFormat: time.RFC3339,
	})
	l.SetOutput(os.Stdout)

	lvl := config.Level
	if lvl == 0 {
		lvl = 6
	}

	l.SetLevel(logrus.Level(lvl))

	return l
}
