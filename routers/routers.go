package routers

import (
	"SimpleGinDemo/controllers"
	"SimpleGinDemo/middlewares"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func InitRouter() {
	gin.DisableConsoleColor()
	f, _ := os.Create("Transfer.log")
	gin.DefaultWriter = io.MultiWriter(f)
	router := gin.Default()
	// 信任代理来源
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return
	}
	// 要在路由组之前全局使用「跨域中间件」, 否则OPTIONS会返回404
	router.Use(middlewares.Cors())
	// 使用 session(cookie-based)
	//router.Use(sessions.Sessions("myyyyysession", Sessions.Store))
	router.StaticFS("/public", http.Dir("./web/static"))
	router.LoadHTMLGlob("web/template/*")
	v1 := router.Group("/")
	{
		v1.POST("/upload", controllers.FileUpload)
		v1.GET("/getFile", controllers.GetFile)
	}

	if err := router.Run("0.0.0.0:8888"); err != nil {
		os.Exit(1)
	}
}
