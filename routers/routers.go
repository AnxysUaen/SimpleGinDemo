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

	router.StaticFS("/static", http.Dir("./@websrc/static"))
	router.LoadHTMLFiles("./web/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./web/favicon.ico")
	})

	fileMgr := router.Group("/fileMgr")
	{
		fileMgr.POST("/getList", controllers.GetList)
		fileMgr.POST("/delete", controllers.FileUpload)
		fileMgr.POST("/upload", controllers.FileUpload)
		fileMgr.GET("/getFile", controllers.GetFile)
	}

	if err := router.Run("0.0.0.0:8888"); err != nil {
		os.Exit(1)
	}
}
