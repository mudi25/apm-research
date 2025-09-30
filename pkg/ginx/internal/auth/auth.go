package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"research-apm/pkg/errors"
	"research-apm/pkg/errors/codes"
	"research-apm/pkg/ginx/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func VerifyHMAC(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		type Header struct {
			Timestamp string `header:"X-Auth-Timestamp" binding:"required"`
			Signature string `header:"X-Auth-Signature" binding:"required"`
		}
		var header Header
		if err := c.ShouldBindHeader(&header); err != nil {
			response.Abort(c, errors.New(codes.Unauthorized, "invalid request auth", err))
			return
		}
		if isExpired(header.Timestamp) {
			response.Abort(c, errors.New(codes.Unauthorized, "invalid request auth", fmt.Errorf("timestamp expired")))
			return
		}

		payload := fmt.Sprintf("%s|%s|%s", c.Request.URL.Path, c.Request.Method, header.Timestamp)
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(payload))

		if !hmac.Equal([]byte(hex.EncodeToString(mac.Sum(nil))), []byte(header.Signature)) {
			response.Abort(c, errors.New(codes.Unauthorized, "invalid request auth", fmt.Errorf("invalid timestamp and signature")))
			return
		}
		c.Next()
	}
}

func isExpired(timestamp string) bool {
	tsInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return true
	}

	t := time.Unix(tsInt, 0)
	now := time.Now()

	// Valid jika t >= now - 2 menit dan t <= now + 1 menit
	if t.After(now.Add(-2*time.Minute)) && t.Before(now.Add(1*time.Minute)) {
		return false // masih valid
	}

	return true // expired
}
