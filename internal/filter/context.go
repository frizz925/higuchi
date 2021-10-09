package filter

import (
	"net"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogFields struct {
	Worker      int
	Proto       string
	Server      string
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
	lf := c.LogFields
	zf := []zapcore.Field{zap.Int("worker", lf.Worker)}
	if lf.Proto != "" {
		zf = append(zf, zap.String("proto", lf.Proto))
	}
	if lf.Server != "" {
		zf = append(zf, zap.String("server", lf.Server))
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
