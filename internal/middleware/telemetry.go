package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/leonardonicola/golerplate/pkg/constants"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func TracingMiddleware() gin.HandlerFunc {
	return otelgin.Middleware(constants.TRACER_NAME)
}
