package routers

import (
	"SimpleGinDemo/controllers"
	"SimpleGinDemo/middlewares"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.DisableConsoleColor()
	f, _ := os.Create("Transfer.log")
	gin.DefaultWriter = io.MultiWriter(f)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// 信任代理来源
	//if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
	//	return
	//}
	// 要在路由组之前全局使用「跨域中间件」, 否则OPTIONS会返回404
	router.Use(middlewares.Cors())

	router.StaticFS("/assets", http.Dir("./web/assets"))
	router.LoadHTMLFiles("./web/index.html")
	router.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/favicon.svg", func(c *gin.Context) {
		c.File("./web/favicon.svg")
	})

	fileMgr := router.Group("/fileMgr")
	{
		fileMgr.POST("/getList", controllers.GetList)
		fileMgr.POST("/delFile", controllers.DelFile)
		fileMgr.POST("/upload", controllers.FileUpload)
		fileMgr.GET("/getFile", controllers.GetFile)
	}
	tools := router.Group("/tools")
	{
		tools.POST("/sendText", controllers.SendText)
		tools.POST("/getText", controllers.GetText)
	}
	mongos := router.Group("/mongos")
	{
		mongos.POST("/saveLog", controllers.SaveLog)
	}

	if err := router.Run("0.0.0.0:8888"); err != nil {
		os.Exit(1)
	}
}
