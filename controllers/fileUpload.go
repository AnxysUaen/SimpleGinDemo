package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FileUpload(c *gin.Context) {
	if file, err := c.FormFile("file"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMsg": "服务文件接收失败!",
		})
	} else {
		// 改为配置
		err := c.SaveUploadedFile(file, "../uploads/"+file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errMsg": "文件保存失败",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"succMsg": "ok",
			})
		}
	}
}
