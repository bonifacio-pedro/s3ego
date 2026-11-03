package rest

import (
	"encoding/xml"
	"net/http"

	"github.com/bonifacio-pedro/s3ego/internal/model"
	"github.com/gin-gonic/gin"
)

// RespondS3Error writes a standardized AWS S3 XML error payload.
func RespondS3Error(c *gin.Context, statusCode int, code, message string) {
	c.Header("Content-Type", "application/xml")
	reqID := c.Writer.Header().Get("x-amz-request-id")
	if reqID == "" {
		reqID = "S3EGO-REQ"
	}

	errResp := model.S3ErrorResponse{
		Code:      code,
		Message:   message,
		Resource:  c.Request.URL.Path,
		RequestID: reqID,
	}

	payload, err := xml.MarshalIndent(errResp, "", "  ")
	if err != nil {
		c.String(http.StatusInternalServerError, `<Error><Code>InternalError</Code></Error>`)
		return
	}

	c.Data(statusCode, "application/xml; charset=utf-8", append([]byte(xml.Header), payload...))
}
