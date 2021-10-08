package filter

import (
	"net"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogFields struct {
	Proto       string
	Listener    string
	Source      string
	Destination string
	User        string
}

type Context struct {
	net.Conn

	Logger    *zap.Logger
	LogFields LogFields

	rootLogger *zap.Logger
}

func NewContext(conn net.Conn, logger *zap.Logger, fields ...LogFields) *Context {
	ctx := &Context{
		Conn:       conn,
		Logger:     logger,
		rootLogger: logger,
	}
	if len(fields) > 0 {
		ctx.LogFields = fields[0]
		ctx.UpdateLogger()
	}
	return ctx
}

func (c *Context) UpdateLogger() {
	zf, lf := []zapcore.Field{}, c.LogFields
	if lf.Proto != "" {
		zf = append(zf, zap.String("proto", lf.Proto))
	}
	if lf.Listener != "" {
		zf = append(zf, zap.String("listener", lf.Listener))
	}
	if lf.Source != "" {
		zf = append(zf, zap.String("src", lf.Source))
	}
	if lf.Destination != "" {
		zf = append(zf, zap.String("dst", lf.Destination))
	}
	if lf.User != "" {
		zf = append(zf, zap.String("user", lf.User))
	}
	c.Logger = c.rootLogger.With(zf...)
}
