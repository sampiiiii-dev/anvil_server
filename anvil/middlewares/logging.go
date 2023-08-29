package middlewares

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ScribeLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			logger.Info("incoming request",
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Int("status", res.Status),
			)
			return next(c)
		}
	}
}
