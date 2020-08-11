// Copyright © 2019 Annchain Authors <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package mylog

import (
	"fmt"
	"github.com/annchain/commongo/utilfuncs"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"path/filepath"
)

type LogConfig struct {
	MaxSize    int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
	LogDir     string
	OutputFile string
}

func RotateLog(abspath string, config LogConfig) *lumberjack.Logger {
	logFile := &lumberjack.Logger{
		Filename:   abspath,
		MaxSize:    config.MaxSize, // megabytes
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAgeDays, //days
		Compress:   config.Compress,   // disabled by default
	}
	return logFile
}

func InitLogger(logger *logrus.Logger, config LogConfig) *logrus.Logger {
	var writer io.Writer
	if config.LogDir != "" {
		folderPath, err := filepath.Abs(config.LogDir)
		utilfuncs.PanicIfError(err, fmt.Sprintf("Error on parsing log path: %s", config.LogDir))

		abspath, err := filepath.Abs(path.Join(config.LogDir, config.OutputFile))
		utilfuncs.PanicIfError(err, fmt.Sprintf("Error on parsing log file path: %s", config.LogDir))

		err = os.MkdirAll(folderPath, os.ModePerm)
		utilfuncs.PanicIfError(err, fmt.Sprintf("Error on creating log dir: %s", folderPath))

		//logFile, err := os.OpenFile(abspath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		utilfuncs.PanicIfError(err, fmt.Sprintf("Error on creating log file: %s", abspath))
		//write  a message to just one  files

		writer = io.MultiWriter(logger.Out, RotateLog(abspath, config))
	} else {
		writer = logger.Out
	}
	newLogger := &logrus.Logger{
		Level:        logger.Level,
		Formatter:    logger.Formatter,
		Out:          writer,
		Hooks:        logger.Hooks,
		ExitFunc:     logger.ExitFunc,
		ReportCaller: logger.ReportCaller,
	}
	return newLogger
}

func LogInit(level logrus.Level) {
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "15:04:05.000000"
	Formatter.FullTimestamp = true
	Formatter.ForceColors = true
	logrus.SetFormatter(Formatter)
	logrus.SetLevel(level)
	//logrus.SetReportCaller(true)
}

func LogLevel(text string) (level logrus.Level) {
	// Only log the warning severity or above.
	switch text {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		fmt.Println("Unknown level", text, "Set to INFO")
		return logrus.InfoLevel
	}

}
