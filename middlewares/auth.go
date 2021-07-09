package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wanhuasong/genericfs/config"
	"github.com/wanhuasong/genericfs/controllers"
	"github.com/wanhuasong/genericfs/utils"
)

const (
	EndpointMkzip     = "/mkzip"
	EndpointPersist   = "/persist"
	EndpointPreupload = "/preupload"
	EndpointDownload  = "/download"
	EndpointUpload    = "/upload"
)

func Auth(c *gin.Context) {
	var token string
	var origin string
	params := make(map[string]string)
	switch c.Request.Method {
	case http.MethodPost:
		switch c.ContentType() {
		case gin.MIMEJSON:
			endpoint := c.Request.URL.Path
			b, err := ReadBody(c)
			if err != nil {
				utils.Response(c, http.StatusBadRequest, fmt.Sprintf("%+v", err))
				c.Abort()
				return
			}
			switch endpoint {
			case EndpointPreupload:
				var req controllers.PreuploadRequest
				if err := json.Unmarshal(b, &req); err != nil {
					utils.Response(c, http.StatusBadRequest, fmt.Sprintf("%+v", err))
					c.Abort()
					return
				}
				log.Printf("Auth preupload req: %+v", req)
				params["t"] = strconv.FormatInt(req.SignTime, 10)
				params["callback_url"] = req.CallbackURL
				params["callback_body"] = req.CallbackBody
				params["uuid"] = req.UUID
				params["reference_type"] = strconv.Itoa(req.ReferenceType)
				params["reference_id"] = req.ReferenceID
				params["team_uuid"] = req.TeamUUID
				params["project_uuid"] = req.ProjectUUID
				params["owner_uuid"] = req.OwnerUUID
				params["modifier"] = req.Modifier
				params["type"] = strconv.Itoa(req.Type)
				params["source"] = strconv.Itoa(req.Source)
				params["ext_id"] = req.ExtID
				params["name"] = req.Name
				params["status"] = strconv.Itoa(req.Status)
				params["create_time"] = strconv.FormatInt(req.CreateTime, 10)
				params["description"] = req.Description
				params["modify_time"] = strconv.FormatInt(req.ModifyTime, 10)
				params["is_public"] = strconv.FormatBool(req.IsPublic)
				token = req.Token

			case EndpointMkzip:
				var req controllers.MkzipRequest
				if err := json.Unmarshal(b, &req); err != nil {
					utils.Response(c, http.StatusBadRequest, fmt.Sprintf("%+v", err))
					c.Abort()
					return
				}
				log.Printf("Auth mkzip req: %+v", req)
				params["callback_url"] = req.CallbackURL
				params["source_hashes"] = req.SourceHashes
				params["target_hash"] = req.TargetHash
				params["t"] = strconv.FormatInt(req.SignTime, 10)
				token = req.Token

			case EndpointPersist:
				var req controllers.PersistRequest
				if err := json.Unmarshal(b, &req); err != nil {
					utils.Response(c, http.StatusBadRequest, fmt.Sprintf("%+v", err))
					c.Abort()
					return
				}
				log.Printf("Auth persist req: %+v", req)
				params["t"] = strconv.FormatInt(req.SignTime, 10)
				params["hash"] = req.Hash
				token = req.Token

			default:
				utils.Response(c, http.StatusNotFound, "Endpoint not found")
				c.Abort()
				return
			}

		case gin.MIMEMultipartPOSTForm:
			var err error
			origin, token, err = utils.ParseUploadToken(c)
			if err != nil {
				utils.Response(c, http.StatusBadRequest, fmt.Sprintf("%+v", err))
				c.Abort()
				return
			}

		default:
			utils.Response(c, http.StatusBadRequest, "Bad request")
			c.Abort()
			return
		}

	case http.MethodGet:
		c.Next()
		return

	default:
		utils.Response(c, http.StatusBadRequest, "Unsupported method")
		c.Abort()
		return
	}

	if origin == "" && len(params) != 0 {
		t, found := params["t"]
		if !found {
			utils.Response(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}
		signTime, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			log.Printf("Parse sign time failed: %+v", err)
			utils.Response(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}
		if config.Cfg.TokenExpire+signTime < time.Now().Unix() {
			utils.Response(c, http.StatusUnauthorized, "Token expired")
			c.Abort()
			return
		}

		origin = utils.SortParams(params)
	}
	if err := utils.VerifySig(origin, token); err != nil {
		utils.Response(c, http.StatusUnauthorized, "Invalid token")
		c.Abort()
		return
	}
	c.Next()
}

func ReadBody(c *gin.Context) ([]byte, error) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return b, nil
}
