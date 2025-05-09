package gologger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	Level_INFO    = "[INF]"
	Level_WARN    = "[WRN]"
	Level_ERROR   = "[ERR]"
	Level_DEBUG   = "[DEB]"
	Level_Success = "[SUC]"
)

type MsgInfo struct {
	Level string
	Msg   string
}

func logFile(ctx context.Context, level string, msg string) {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		runtime.EventsEmit(ctx, "gologger", &MsgInfo{
			Level: Level_ERROR,
			Msg:   Msg(err.Error()),
		})
	}
	line := fmt.Sprintf("%v %v\n", level, msg)
	_, err = file.Write([]byte(line))
	if err != nil {
		runtime.EventsEmit(ctx, "gologger", &MsgInfo{
			Level: Level_ERROR,
			Msg:   Msg(err.Error()),
		})
	}
}
func Info(ctx context.Context, i interface{}) {
	logFile(ctx, Level_INFO, Msg(i))
	runtime.EventsEmit(ctx, "gologger", &MsgInfo{
		Level: Level_INFO,
		Msg:   Msg(i),
	})
}

func Warning(ctx context.Context, i interface{}) {
	logFile(ctx, Level_INFO, Msg(i))
	runtime.EventsEmit(ctx, "gologger", &MsgInfo{
		Level: Level_WARN,
		Msg:   Msg(i),
	})
}

func Error(ctx context.Context, i interface{}) {
	logFile(ctx, Level_INFO, Msg(i))
	runtime.EventsEmit(ctx, "gologger", &MsgInfo{
		Level: Level_ERROR,
		Msg:   Msg(i),
	})
}

func Debug(ctx context.Context, i interface{}) {
	logFile(ctx, Level_INFO, Msg(i))
	runtime.EventsEmit(ctx, "gologger", &MsgInfo{
		Level: Level_DEBUG,
		Msg:   Msg(i),
	})
}

func Success(ctx context.Context, i interface{}) {
	logFile(ctx, Level_INFO, Msg(i))
	runtime.EventsEmit(ctx, "gologger", &MsgInfo{
		Level: Level_Success,
		Msg:   Msg(i),
	})
}

func Msg(i interface{}) string {
	return fmt.Sprintf("[%s] %v", currentTime(), i)
}

func currentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 当出现大量错误时，每个5s打印一次错误日志，偷懒式写法
func IntervalError(ctx context.Context, i interface{}) {
	currentTime := time.Now()
	unixTimestamp := currentTime.Unix()
	if unixTimestamp%5 == 0 {
		Error(ctx, i)
	}
}
