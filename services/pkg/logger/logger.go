package logger

import (
	"os"
	"github.com/sirupsen/logrus"
	formatter "github.com/antonfisher/nested-logrus-formatter"
)

var Logger logrus.Logger

func Init() {
	Logger = logrus.Logger{
		Out:   os.Stderr,
        Level: logrus.DebugLevel,
        Formatter: &formatter.Formatter{
            TimestampFormat: "[2006-01-02 15:04:05]",
            HideKeys: true,
			NoColors: false,
        },
	}
}