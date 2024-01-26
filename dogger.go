package dogger

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once
var log zerolog.Logger

var correlationKey string

func Get() zerolog.Logger {
	once.Do(func() {

		correlationKey = os.Getenv("CORRELATION_KEY")
		if correlationKey == "" {
			correlationKey = "corrId"
		}
		logLevel := os.Getenv("LOG_LEVEL")
		serviceName := os.Getenv("SERVICE_NAME")
		serviceEnv := os.Getenv("SERVICE_ENV")
		logMaxSize, err := strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
		if err != nil {
			logMaxSize = 100 // default to 100MB
		}
		logMaxBackups, err := strconv.Atoi(os.Getenv("LOG_MAX_BACKUPS"))
		if err != nil {
			logMaxBackups = 0 // default to 0 backups
		}
		logMaxAge, err := strconv.Atoi(os.Getenv("LOG_MAX_AGE"))
		if err != nil {
			logMaxAge = 0 // default to 0 days
		}

		// If logMaxBackups and logMaxAge are both 0, no old log files will be deleted.

		logCompressed, err := strconv.ParseBool(os.Getenv("LOG_COMPRESSED"))
		if err != nil {
			logCompressed = true // default to compressed
		}

		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		intLogLevel, err := strconv.Atoi(logLevel)
		if err != nil {
			intLogLevel = int(zerolog.DebugLevel) // default to DEBUG
		}

		fileLogger := &lumberjack.Logger{
			Filename:   fmt.Sprintf("%s.log", serviceName),
			MaxSize:    logMaxSize,
			MaxBackups: logMaxBackups,
			MaxAge:     logMaxAge,
			Compress:   logCompressed,
		}

		var output io.Writer
		if serviceEnv != "development" {
			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		} else {
			// CAUTION: zerolog.ConsoleWriter is not safe for concurrent use.
			// Use it ONLY in DEVELOPMENT.
			output = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			}
			// saving into stderr, file and stdout FOR DEVELOPMENT ONLY
			output = zerolog.MultiLevelWriter(fileLogger, output)
		}

		var gitRevision string
		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		log = zerolog.New(output).
			Level(zerolog.Level(intLogLevel)).
			With().Timestamp().Str("app", serviceName).Logger()

		if v, ok := os.LookupEnv("USING_GIT_REVISION"); ok && v == "true" {
			log = log.With().Str("git_revision", gitRevision).Logger()
		}

		if v, ok := os.LookupEnv("USING_GO_VERSION"); ok && v == "true" {
			log = log.With().Str("go_version", buildInfo.GoVersion).Logger()
		}

		if v, ok := os.LookupEnv("USING_PID"); ok && v == "true" {
			log = log.With().Str("pid", fmt.Sprint(os.Getpid())).Logger()
		}
	})

	return log
}

func mainLogger(level string, corrId string, event string, dataObj interface{}, err error) {
	// Buffered Channel of type Boolean
	done := make(chan bool, 1)

	go func(done chan bool, level string) {
		defer func() {
			// Sending value to channel
			done <- true
			return
		}()

		data := map[string]any{"data": dataObj}

		l := Get()

		var logLevelFunc func() *zerolog.Event
		switch level {
		case "info":
			logLevelFunc = l.Info
		case "warn":
			logLevelFunc = l.Warn
		case "error":
			logLevelFunc = l.Error().Stack
		case "debug":
			logLevelFunc = l.Debug
		}

		logLevelFunc().
			Str(correlationKey, corrId).
			Str("event", event).Fields(data).Err(err).
			Send()
		return
	}(done, level)

	// Waiting to receive value from channel
	<-done
	return
}

func LogInfo(corrId string, event string, dataObj interface{}) {
	mainLogger("info", corrId, event, dataObj, nil)
}

func LogError(corrId string, event string, dataObj interface{}, err error) {
	mainLogger("error", corrId, event, dataObj, err)
}

func LogDebug(corrId string, event string, dataObj interface{}) {
	mainLogger("debug", corrId, event, dataObj, nil)
}

func LogFatal(corrId string, event string, dataObj interface{}) {
	mainLogger("fatal", corrId, event, dataObj, nil)
}
