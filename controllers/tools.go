package controllers

import (
	"github.com/atotto/clipboard"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendText(c *gin.Context) {
	saveText := c.DefaultPostForm("text", "")
	err := clipboard.WriteAll(saveText)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errMsg": "发送失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"succMsg": "ok",
	})
}

func GetText(c *gin.Context) {
	content, err := clipboard.ReadAll()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errMsg": "获取失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"succMsg": "ok",
		"data":    content,
	})
}
