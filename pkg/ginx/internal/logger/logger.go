package logger

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type Logging struct {
	TraceID           string         `json:"traceId"`
	AppName           string         `json:"appName"`
	Method            string         `json:"method"`
	Path              string         `json:"path"`
	ElapsedTime       int64          `json:"elapsedTime"`
	ClientIP          string         `json:"clientIp"`
	Site              string         `json:"site"`
	Environment       string         `json:"environment"`
	ApkVersion        string         `json:"apkVersion"`
	DBApkVersion      string         `json:"dbApkVersion"`
	RequestUser       map[string]any `json:"requestUser"`
	RequestQuery      map[string]any `json:"requestQuery"`
	RequestBody       map[string]any `json:"requestBody"`
	ResponseCode      int            `json:"responseCode"`
	ResponseBody      map[string]any `json:"responseBody"`
	AdditionalContent map[string]any `json:"additionalContent"`
	Timestamp         time.Time      `json:"timestamp"`
}

func NewLogger(
	appName string,
	appSite string,
	appEnv string,
	appVersion string,
	appDbVersion string,
	logCh chan Logging,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		timestamp := time.Now()

		// Capture response body using a tee writer
		buf := new(bytes.Buffer)
		tee := &teeResponseWriter{ctx.Writer, buf}
		ctx.Writer = tee

		// Read and clone request body if it's JSON
		var requestBody map[string]any = nil
		if ctx.Request.Header.Get("Content-Type") == "application/json" {
			if err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody); err != nil {
				requestBody = nil
			}
			jbody, _ := json.Marshal(&requestBody)
			ctx.Request.Body = io.NopCloser(bytes.NewReader(jbody)) // restore body
		}

		// Proceed with request
		ctx.Next()

		// Parse JSON response
		var responseBody map[string]any = nil
		json.Unmarshal(buf.Bytes(), &responseBody)

		// Decode auth user from custom header

		var reqUser map[string]any = nil
		if reqUserHeader := ctx.GetHeader("X-Auth-User"); reqUserHeader != "" {
			if b, err := base64.StdEncoding.DecodeString(reqUserHeader); err == nil {
				json.Unmarshal(b, &reqUser)
			}
		}

		// Convert query params into map[string]any
		var requestQuery map[string]any = nil
		if rawQuery := ctx.Request.URL.Query(); len(rawQuery) > 0 {
			requestQuery = make(map[string]any)
			for k, v := range rawQuery {
				if len(v) == 1 {
					requestQuery[k] = v[0]
				} else {
					requestQuery[k] = v
				}
			}
		}
		logCh <- Logging{
			TraceID:           ctx.GetHeader("X-Trace-ID"),
			AppName:           appName,
			Method:            ctx.Request.Method,
			Path:              ctx.Request.URL.Path,
			ElapsedTime:       time.Since(timestamp).Milliseconds(),
			ClientIP:          ctx.ClientIP(),
			Site:              appSite,
			Environment:       appEnv,
			ApkVersion:        appVersion,
			DBApkVersion:      appDbVersion,
			RequestUser:       reqUser,
			RequestQuery:      requestQuery,
			RequestBody:       requestBody,
			ResponseCode:      ctx.Writer.Status(),
			ResponseBody:      responseBody,
			AdditionalContent: nil,
			Timestamp:         timestamp,
		}
	}
}

// teeResponseWriter is a wrapper around gin.ResponseWriter
// that duplicates writes to an internal buffer so the response body can be logged.
type teeResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write writes the data to both the response writer and the internal buffer.
func (w *teeResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// Status returns the HTTP status code of the response.
func (w *teeResponseWriter) Status() int {
	return w.ResponseWriter.Status()
}
