package controllers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wanhuasong/genericfs/config"
	"github.com/wanhuasong/genericfs/models"
	"github.com/wanhuasong/genericfs/utils"
)

type PreuploadRequest struct {
	SignTime int64  `json:"t"` // time.second
	Token    string `json:"token"`

	*models.Resource
}

type MkzipRequest struct {
	SignTime     int64  `json:"t"`
	Token        string `json:"token"`
	CallbackURL  string `json:"callback_url"`
	SourceHashes string `json:"source_hashes"`
	TargetHash   string `json:"target_hash"`
}

type PersistRequest struct {
	Token    string `json:"token"`
	SignTime int64  `json:"t"`
	Hash     string `json:"hash"`
}

func Download(c *gin.Context) {
	hash := c.Param("hash")
	if len(hash) == 0 {
		utils.Response(c, http.StatusBadRequest, "Bad request")
		return
	}

	resource, err := models.GetResourceByHash(hash)
	if err != nil {
		utils.Response(c, http.StatusBadRequest, fmt.Sprintf("Get resource failed: %+v", err))
		return
	}

	rawQuery := c.Request.URL.RawQuery
	rawQuery, err = url.QueryUnescape(rawQuery)
	if err != nil {
		utils.Response(c, http.StatusBadRequest, fmt.Sprintf("Unescape query failed: %+v", err))
		return
	}
	log.Printf("Download raw query: %s", rawQuery)

	op := strings.Split(rawQuery, "&")[0]
	if strings.IndexByte(op, '=') != -1 {
		op = ""
	}
	if op != "" {
		handleOP(op)
	}

	if !resource.IsPublic {
		err = authDownload(hash, rawQuery, op)
		if err != nil {
			utils.Response(c, http.StatusUnauthorized, "Invalid token")
			return
		}
	}

	isPreview := true
	filename := c.Query("attname")
	if filename == "" {
		filename = resource.Name
		isPreview = false
	}

	fp := filepath.Join(config.Cfg.StoreDir, hash)
	b, err := os.ReadFile(fp)
	if err != nil {
		utils.Response(c, http.StatusNotFound, "File not found")
		return
	}

	contentType, err := utils.GetFileContentType(fp)
	if err != nil {
		utils.Response(c, http.StatusNotFound, "File not found")
		return
	}

	var contentDisposition string
	if isPreview {
		contentDisposition = fmt.Sprintf("inline; filename=%s", filename)
	} else {
		contentDisposition = fmt.Sprintf("attachment; filename=%s", filename)
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", contentDisposition)
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(b)))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Writer.Write(b)
}

func handleOP(op string) {
	log.Printf("Handle op: %s", op)
}

func authDownload(hash, query, op string) error {
	items := strings.Split(query, "&")
	var token string
	params := map[string]string{
		"hash": hash,
	}
	if op != "" {
		params["op"] = op
	}
	for _, item := range items {
		idx := strings.IndexByte(item, '=')
		if idx == -1 || len(item) <= idx+1 {
			continue
		}

		key := item[:idx]
		if key == "attname" {
			continue
		}

		val := item[idx+1:]
		if key == utils.SigKey {
			token = val
		} else {
			params[key] = val
		}

		if key == "e" {
			expire, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			if time.Now().Unix() > expire {
				return errors.New("Token expired")
			}
		}
	}
	log.Printf("Download params: %+v", params)
	origin := utils.SortParams(params)
	return utils.VerifySig(origin, token)
}

func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.Response(c, http.StatusBadRequest, fmt.Sprintf("Invalid form file: %+v", err))
		return
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		utils.Response(c, http.StatusBadRequest, fmt.Sprintf("Read file failed: %+v", err))
		return
	}

	// fb := make([]byte, len(content))
	// copy(fb, content)
	hash, err := utils.GetEtag(content)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, fmt.Sprintf("Get etag failed: %+v", err))
		return
	}

	fileSize := len(content)
	mime := header.Header.Get("Content-Type")
	savePath := filepath.Join(config.Cfg.StoreDir, hash)
	f, err := os.Create(savePath)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, fmt.Sprintf("Create file failed: %+v", err))
		return
	}
	f.Write(content)
	f.Close()

	filename := header.Filename
	downloadURL, err := updateResource(c, filename, hash, mime, fileSize)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, fmt.Sprintf("Update resource failed: %+v", err))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success",
		"url":     downloadURL,
		"version": 1,
		"name":    filename,
		"hash":    hash,
		"mime":    mime,
		"size":    fileSize,
	})
}

func updateResource(c *gin.Context, filename, hash, mime string, fileSize int) (string, error) {
	origin, _, _ := utils.ParseUploadToken(c)
	var uuid string
	items := strings.Split(origin, "&")
	for _, item := range items {
		kv := strings.Split(item, "=")
		if len(kv) != 2 {
			continue
		}
		if kv[0] == "uuid" {
			uuid = kv[1]
		}
	}

	// Generate download url while resource is public
	var downloadURL string
	if uuid == "" {
		resource := &models.Resource{
			IsPublic: true,
			ExtID:    hash,
			Name:     filename,
		}
		err := models.AddResource(resource)
		if err != nil {
			return "", err
		}
		downloadURL = fmt.Sprintf("%s/download/%s", config.Cfg.BaseURL, hash)
	} else {
		resource, err := models.GetResourceByUUID(uuid)
		if err != nil {
			return "", err
		}

		if resource.IsPublic {
			downloadURL = fmt.Sprintf("%s/download/%s", config.Cfg.BaseURL, hash)
		}

		resource.Name = filename
		resource.ExtID = hash
		resource.ModifyTime = time.Now().UnixNano() / int64(time.Microsecond)
		defer uploadCallback(resource, mime, fileSize)

		err = models.UpdateResource(resource)
		if err != nil {
			return "", err
		}
	}
	return downloadURL, nil
}

func uploadCallback(resource *models.Resource, mime string, fileSize int) {
	resource.CallbackBody = fillCallbackBody(resource.CallbackBody, resource, mime, fileSize)

	log.Printf("Upload callback url: %s", resource.CallbackURL)
	log.Printf("Upload callback body: %s", resource.CallbackBody)
	payload := strings.NewReader(resource.CallbackBody)
	req, err := http.NewRequest(http.MethodPost, resource.CallbackURL, payload)
	if err != nil {
		log.Printf("New upload callback request failed: %+v", err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Upload callback failed: %+v", err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Read upload callback response body failed: %+v", err)
		return
	}
	log.Printf("Upload callback response body: %s", string(body))
}

func fillCallbackBody(callbackBody string, resource *models.Resource, mime string, fileSize int) string {
	params := make(map[string]string)
	items := strings.Split(callbackBody, "&")
	re := regexp.MustCompile(`^\$(.+)$`)

	var w int
	var h int
	if strings.HasPrefix(mime, "image/") {
		p := filepath.Join(config.Cfg.StoreDir, resource.ExtID)
		w, h = getImageDimension(p)
	}

	for _, item := range items {
		kv := strings.Split(item, "=")
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		val := kv[1]

		// process $() val
		if re.MatchString(val) {
			switch val[2 : len(val)-1] {
			case "key":
				val = resource.ExtID
			case "name":
				val = resource.Name
			case "size":
				val = strconv.Itoa(fileSize)
			case "suffix":
				dotIdx := strings.LastIndex(resource.Name, ".")
				if dotIdx == -1 {
					val = ""
				} else {
					val = resource.Name[dotIdx+1:]
				}
			case "mimeType":
				val = mime
			case "exif":
				val = ""
			case "imageInfo.width":
				val = strconv.Itoa(w)
			case "imageInfo.height":
				val = strconv.Itoa(h)
			default:
				log.Printf("Unknown upload callback val: %s", val)
				continue
			}
		}
		params[key] = val
	}

	callbackBody = ""
	for k, v := range params {
		if callbackBody != "" {
			callbackBody += "&"
		}
		callbackBody += fmt.Sprintf("%s=%s", k, v)
	}
	return callbackBody
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("Open image file failed: %+v", err)
		return 0, 0
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("Decode image config failed: %+v", err)
		return 0, 0
	}
	return image.Width, image.Height
}

func Preupload(c *gin.Context) {
	var req PreuploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "Bad request")
		return
	}
	log.Printf("Preupload req: %+v", req)

	if err := models.AddResource(req.Resource); err != nil {
		utils.Response(c, http.StatusInternalServerError, fmt.Sprintf("Failed to add resource: %+v", err))
		return
	}
	utils.Response(c, http.StatusOK, "Success")
}

func Mkzip(c *gin.Context) {
	var req MkzipRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Response(c, http.StatusBadRequest, "Bad request")
		return
	}
	log.Printf("Mkzip req: %+v", req)

	hashes := strings.Split(req.SourceHashes, ",")
	err = zipFiles(req.TargetHash, hashes)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, fmt.Sprintf("Zip files failed: %+v", err))
		return
	}

	filename := fmt.Sprintf("%s.zip", req.TargetHash)
	resource := &models.Resource{
		ExtID: req.TargetHash,
		Name:  filename,
	}
	err = models.AddResource(resource)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, fmt.Sprintf("Add resource failed: %+v", err))
		return
	}

	defer mkzipCallback(req.CallbackURL)
	utils.Response(c, http.StatusOK, "Success")
}

func zipFiles(hash string, hashes []string) error {
	fp := filepath.Join(config.Cfg.StoreDir, hash)
	targetFile, err := os.Create(fp)
	if err != nil {
		return err
	}

	// Create a new zip archive.
	w := zip.NewWriter(targetFile)

	for _, h := range hashes {
		resource, err := models.GetResourceByHash(h)
		if err != nil {
			log.Printf("Get resource failed: %+v", err)
			continue
		}

		p := filepath.Join(config.Cfg.StoreDir, h)
		pf, err := os.Open(p)
		if err != nil {
			log.Printf("Open file failed: %+v", err)
			continue
		}

		f, err := w.Create(resource.Name)
		if err != nil {
			log.Printf("Create file failed: %+v", err)
			continue
		}

		_, err = io.Copy(f, pf)
		if err != nil {
			log.Printf("Copy file content failed: %+v", err)
			continue
		}

		pf.Close()
	}

	// Make sure to check the error on Close.
	return w.Close()
}

func mkzipCallback(callbackURL string) {
	log.Printf("Mkzip callback url: %s", callbackURL)

	data := map[string]interface{}{
		"code": 0,
	}
	b, _ := json.Marshal(data)
	r := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodPost, callbackURL, r)
	if err != nil {
		log.Printf("New mkzip callback request failed: %+v", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Mkzip callback failed: %+v", err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Read mkzip callback response body failed: %+v", err)
		return
	}
	log.Printf("Mkzip callback response body: %s", string(body))
}

func Persist(c *gin.Context) {
	var req PersistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, "Bad request")
		return
	}
	log.Printf("Persist req: %+v", req)

	utils.Response(c, http.StatusOK, "Success")
}
