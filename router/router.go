package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wanhuasong/genericfs/fs"
	"github.com/wanhuasong/genericfs/utils"
)

const SigKey = "token"

func Run() error {
	api := gin.Default()

	api.Use(auth)

	api.GET("/download", fs.Download)
	api.POST("/upload", fs.Upload)
	api.POST("/preupload", fs.Preupload)
	api.POST("/mkzip", fs.Mkzip)
	api.POST("/persist", fs.Persist)
	return api.Run(":8080")
}

func auth(c *gin.Context) {
	var sig string
	params := make(map[string]string)
	switch c.Request.Method {
	case http.MethodPost:
		// sign data inside the requestBody
		switch c.ContentType() {
		case gin.MIMEJSON:
			b, _ := ioutil.ReadAll(c.Request.Body)
			data := make(map[string]interface{})
			err := json.Unmarshal(b, &data)
			if err != nil {
				c.JSON(http.StatusBadRequest, map[string]interface{}{
					"code":    http.StatusBadRequest,
					"message": fmt.Sprintf("%+v", err),
				})
				c.Abort()
				return
			}
			log.Printf("data: %+v", data)
		case gin.MIMEMultipartPOSTForm:
			//
		default:
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "Invalid request",
			})
			c.Abort()
			return
		}

	case http.MethodGet:
		// sign data inside the query
		rawQuery := c.Request.URL.RawQuery
		for _, s := range strings.Split(rawQuery, "&") {
			kv := strings.Split(s, "=")
			k := kv[0]
			if k == SigKey {
				if len(kv) != 2 {
					c.JSON(http.StatusBadRequest, map[string]interface{}{
						"code":    http.StatusBadRequest,
						"message": "Bad request",
					})
					c.Abort()
					return
				}
				sig = kv[1]
			} else {
				if len(kv) >= 2 {
					params[k] = kv[1]
				} else {
					params[k] = ""
				}
			}
		}
	}

	origin := sortParams(params)
	err := utils.VerifySig(origin, sig)
	if err != nil {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "Invalid token",
		})
		c.Abort()
		return
	}
	c.Next()
}

func sortParams(params map[string]string) string {
	keys := make([]string, 0)
	for key := range params {
		keys = append(keys, key)
	}
	var s string
	sort.Strings(keys)
	for _, key := range keys {
		if s != "" {
			s += "&"
		}
		s += fmt.Sprintf("%s=%s", key, params[key])
	}
	return s
}
