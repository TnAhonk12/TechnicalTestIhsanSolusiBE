package logger

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

func NewGormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}
