package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(path string) *zap.Logger {
	encoder := getEncoder()

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	infoWriter := getWriteSyncer(filepath.Join(path, "info"))
	warnWriter := getWriteSyncer(filepath.Join(path, "warn"))
	errorWriter := getWriteSyncer(filepath.Join(path, "error"))
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoWriter, infoLevel),
		zapcore.NewCore(encoder, warnWriter, warnLevel),
		zapcore.NewCore(encoder, errorWriter, errorLevel),
	)
	// core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	return logger
}

func getEncoder() zapcore.Encoder {
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	conf.EncodeLevel = zapcore.CapitalLevelEncoder
	conf.EncodeDuration = zapcore.SecondsDurationEncoder
	conf.EncodeCaller = zapcore.ShortCallerEncoder
	conf.TimeKey = "time"
	return zapcore.NewConsoleEncoder(conf)
	// return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getWriteSyncer(path string) zapcore.WriteSyncer {
	workDir, _ := os.Getwd()
	p := filepath.Join(workDir, path)
	writer := getWriter(p)
	w := io.MultiWriter(os.Stdout, writer)
	return zapcore.AddSync(w)
}

func getWriter(filename string) io.Writer {
	hook, _ := rotatelogs.New(
		filename+"_%Y%m%d.log",
		rotatelogs.WithLinkName(filename+".log"),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	return hook
}
