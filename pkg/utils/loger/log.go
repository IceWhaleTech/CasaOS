/*
 * @Author: LinkLeong link@icewhale.com
 * @Date: 2022-06-02 15:09:38
 * @LastEditors: LinkLeong
 * @LastEditTime: 2022-06-16 18:15:58
 * @FilePath: /CasaOS/pkg/utils/loger/log.go
 * @Description:
 * @Website: https://www.casaos.io
 * Copyright (c) 2022 by icewhale, All Rights Reserved.
 */
package loger

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/IceWhaleTech/CasaOS/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var loggers *zap.Logger

func getFileLogWriter() (writeSyncer zapcore.WriteSyncer) {
	// 使用 lumberjack 实现 log rotate
	lumberJackLogger := &lumberjack.Logger{
		Filename: filepath.Join(config.AppInfo.LogPath, fmt.Sprintf("%s.%s",
			config.AppInfo.LogSaveName,
			config.AppInfo.LogFileExt,
		)),
		MaxSize:    100,
		MaxBackups: 60,
		MaxAge:     1,
		Compress:   true,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func LogInit() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	fileWriteSyncer := getFileLogWriter()
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	)
	loggers = zap.New(core)

}

func Info(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	loggers.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	loggers.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	loggers.Error(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	loggers.Warn(message, fields...)
}

func getCallerInfoForLog() (callerFields []zap.Field) {

	pc, file, line, ok := runtime.Caller(2) // 回溯两层，拿到写日志的调用方的函数信息
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName) //Base函数返回路径的最后一个元素，只保留函数名

	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
	return
}
