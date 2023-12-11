package controllers

import (
	"SimpleGinDemo/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
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
		err := c.SaveUploadedFile(file, filepath.Join(savePath, file.Filename))
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
	if _, err := os.Stat(filepath.Clean(fileName)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errMsg": err,
		})
		return
	}
	c.File(filepath.Clean(fileName))
}

func GetList(c *gin.Context) {
	path := c.DefaultPostForm("path", "/")
	info, _ := os.Stat(filepath.Clean(path))
	if info.IsDir() {
		files, _ := os.ReadDir(filepath.Clean(path))
		fileDatas := make([]models.FileData, 0)
		for f := 0; f < len(files); f++ {
			file := files[f]
			fileInfo, _ := file.Info()
			fileDatas = append(fileDatas, models.FileData{
				Name:  file.Name(),
				IsDir: file.IsDir(),
				Size:  fileInfo.Size(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"succMsg": "ok",
			"data":    fileDatas,
		})
	}
}
