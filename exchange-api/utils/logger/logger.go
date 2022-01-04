package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type (
	LoggerInterface interface {
		Debug(args map[string]interface{})
		Error(args map[string]interface{})
		Fatal(args map[string]interface{})
		Info(args map[string]interface{})
		Warn(args map[string]interface{})
	}

	Logger struct {
		Engine *logrus.Entry
	}
)

var LOG_LEVEL = os.Getenv("LOG_LEVEL")

var defaultFields = logrus.Fields{
	"app_name": os.Getenv("APP_NAME"),
	"app_type": os.Getenv("APP_TYPE"),
	"env":      os.Getenv("ENV"),
}

func NewLogger() LoggerInterface {
	engine := logrus.New()

	engine.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	})

	level, _ := logrus.ParseLevel(LOG_LEVEL)

	engine.SetLevel(level)

	logger := engine.WithFields(defaultFields)

	return Logger{
		Engine: logger,
	}
}

func (l Logger) Debug(args map[string]interface{}) {
	message := args["message"]
	delete(args, "message")
	l.Engine.WithFields(args).Debug(message)
}

func (l Logger) Error(args map[string]interface{}) {
	message := args["message"]
	delete(args, "message")
	l.Engine.WithFields(args).Error(message)
}

func (l Logger) Fatal(args map[string]interface{}) {
	message := args["message"]
	delete(args, "message")
	l.Engine.WithFields(args).Fatal(message)
}

func (l Logger) Info(args map[string]interface{}) {
	message := args["message"]
	delete(args, "message")
	l.Engine.WithFields(args).Info(message)
}

func (l Logger) Warn(args map[string]interface{}) {
	message := args["message"]
	delete(args, "message")
	l.Engine.WithFields(args).Warn(message)
}
