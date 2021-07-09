package utils

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	origin := parts[0]
	sig := parts[1]
	// base64 decode
	b, err := base64.StdEncoding.DecodeString(origin)
	if err != nil {
		log.Printf("Decode origin failed: %+v", err)
		return "", "", errors.New("Invalid token")
	}
	// 原加签字符串
	origin = string(b)
	log.Printf("Upload origin: %s", origin)
	// Check token expire
	parts = strings.Split(origin, "&")
	for _, p := range parts {
		kv := strings.Split(p, "=")
		if len(kv) != 2 {
			continue
		}
		if kv[0] == "e" {
			expire, err := strconv.ParseInt(kv[1], 10, 64)
			if err != nil {
				log.Printf("Parse expire failed: %+v", err)
				return "", "", err
			}
			if time.Now().Unix() > expire {
				return "", "", errors.New("Token expired")
			}
		}
	}
	return origin, sig, nil
}
