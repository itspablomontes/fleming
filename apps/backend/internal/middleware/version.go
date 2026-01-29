package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/itspablomontes/fleming/pkg/protocol"
)

// VersionMiddleware adds protocol version headers to HTTP responses.
func VersionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add protocol version header
		c.Header("X-Protocol-Version", protocol.VersionString)
		c.Header("X-Protocol-Schema", protocol.SchemaVersionTimeline) // Default schema

		// Continue with request
		c.Next()
	}
}

// VersionResponse wraps a response with version information.
type VersionResponse struct {
	Version string `json:"version"`
	Schema  string `json:"schema,omitempty"`
	Data    any    `json:"data"`
}

// WrapVersion wraps a response with version metadata.
func WrapVersion(data any, schema string) VersionResponse {
	return VersionResponse{
		Version: protocol.VersionString,
		Schema:  schema,
		Data:    data,
	}
}
