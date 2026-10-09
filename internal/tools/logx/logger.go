package logx

import (
	"context"
	"log"
	"os"
	"strings"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

// Info logs a message at INFO level with the context fields.
func Info(ctx context.Context, msg string) {
	logWithFields("INFO", ctx, msg)
}

// Warn logs a message at WARN level with the context fields.
func Warn(ctx context.Context, msg string) {
	logWithFields("WARN", ctx, msg)
}

// Error logs an error at ERROR level with the context fields.
func Error(ctx context.Context, err error) {
	logWithFields("ERROR", ctx, err.Error())
}

// logWithFields writes one line with the level, the ordered context fields and
// the message. The calldepth (3) makes Lshortfile point at the real caller.
func logWithFields(level string, ctx context.Context, msg string) {
	var b strings.Builder
	b.WriteString(level)
	b.WriteString(" ")
	CtxFields(ctx).Each(func(k, v string) {
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(v)
		b.WriteString(";")
	})
	b.WriteString(" ")
	b.WriteString(msg)
	_ = logger.Output(3, b.String())
}
