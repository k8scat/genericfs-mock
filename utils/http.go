package utils

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const SigKey = "token"

func Response(c *gin.Context, statusCode int, message string) {
	code := statusCode
	if statusCode != http.StatusOK {
		code = http.StatusInternalServerError
	}
	log.Printf("Response %d: %s", statusCode, message)
	c.JSON(code, map[string]interface{}{
		"code":    statusCode,
		"message": message,
	})
}

func ParseUploadToken(c *gin.Context) (string, string, error) {
	token := c.Request.FormValue(SigKey)
	log.Printf("Upload token: %s", token)
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return "", "", errors.New("Invalid token")
	}
	// base64 编码的原加签字符串
	origin := parts[0]
	b, err := base64.StdEncoding.DecodeString(origin)
	if err != nil {
		log.Printf("Failed to decode: %+v", err)
		return "", "", errors.New("Invalid token")
	}
	// 原加签字符串
	origin = string(b)
	// 签名
	sig := parts[1]
	return origin, sig, nil
}
