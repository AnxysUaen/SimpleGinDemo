package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
)

func FileUpload(c *gin.Context) {
	savePath := c.DefaultPostForm("path", "./")
	if info, err := os.Stat(savePath); err != nil || !info.IsDir() {
		c.JSON(http.StatusOK, gin.H{
			"errMsg": "路径无效",
		})
		return
	}
	if file, err := c.FormFile("file"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMsg": "服务文件接收失败!",
		})
	} else {
		err := c.SaveUploadedFile(file, path.Join(savePath, file.Filename))
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

func GetFile(c *gin.Context) {
	fileName, has := c.GetQuery("fileName")
	if !has {
		c.JSON(http.StatusOK, gin.H{
			"errMsg": "缺少字段[fileName]",
		})
		return
	}
	if _, err := os.Stat(path.Join("./", fileName)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errMsg": "文件不存在",
		})
		return
	}
	c.File(path.Join("./", fileName))
}
