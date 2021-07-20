package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wanhuasong/genericfs/controllers"
	"github.com/wanhuasong/genericfs/middlewares"
)

func Run() error {
	api := gin.Default()

	// api.Use(cors.New(cors.Config{
	// 	AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	// 	AllowHeaders: []string{"Content-Type", "X-CSRF-TOKEN", "Accept", "Referer",
	// 		"User-Agent", "Content-Length", "Connection", "Accept-Encoding", "Accept-Language",
	// 		"Cache-Control", "Host", "Origin"},
	// 	ExposeHeaders:   []string{"Content-Type", "Content-Length", "Content-Disposition"},
	// 	AllowAllOrigins: true,
	// 	MaxAge:          12 * time.Hour,
	// }))

	// api.Use(gzip.Gzip(gzip.DefaultCompression))
	api.Use(middlewares.Auth)

	api.GET("/download/:hash", controllers.Download)
	api.POST("/upload", controllers.Upload)
	api.POST("/preupload", controllers.Preupload)
	api.POST("/mkzip", controllers.Mkzip)
	api.POST("/persist", controllers.Persist)
	return api.Run(":14500")
}
